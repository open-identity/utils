package sqlcon

import (
	"database/sql"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	DriverPostgreSQL = "postgres"
	DriverMySQL      = "mysql"
)

// SQLConnection represents a connection to a SQL database.
type SQLConnection struct {
	db  *sqlx.DB
	URL *url.URL
	L   logrus.FieldLogger
	options
}

// NewSQLConnection returns a new SQLConnection.
func NewSQLConnection(db string, l logrus.FieldLogger, opts ...OptionModifier) (*SQLConnection, error) {
	u, err := url.Parse(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if l == nil {
		logger := logrus.New()

		// Basically avoids any logging because no one uses panics
		logger.Level = logrus.PanicLevel

		l = logger
	}

	connection := &SQLConnection{
		URL: u,
		L:   l,
	}

	for _, opt := range opts {
		opt(&connection.options)
	}

	return connection, nil
}

func cleanURLQuery(c *url.URL) *url.URL {
	cleanurl := new(url.URL)
	*cleanurl = *c

	q := cleanurl.Query()
	q.Del("max_conns")
	q.Del("max_idle_conns")
	q.Del("max_conn_lifetime")
	q.Del("parseTime")

	cleanurl.RawQuery = q.Encode()
	return cleanurl
}

// GetDatabaseRetry tries to connect to a database and fails after failAfter.
func (c *SQLConnection) GetDatabaseRetry(maxWait time.Duration, failAfter time.Duration) (*sqlx.DB, error) {
	backOff := backoff.NewExponentialBackOff()
	backOff.MaxInterval = maxWait
	backOff.MaxElapsedTime = failAfter

	if err := backoff.Retry(func() (err error) {
		c.db, err = c.GetDatabase()
		if err != nil {
			return err
		}
		return nil
	}, backOff); err != nil {
		return nil, errors.WithStack(err)
	}

	return c.db, nil
}

// GetDatabase retrusn a database instance.
func (c *SQLConnection) GetDatabase() (*sqlx.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	var err error
	var registeredDriver string

	clean := cleanURLQuery(c.URL)
	if registeredDriver, err = c.registerDriver(); err != nil {
		return nil, errors.Wrap(err, "could not register driver")
	}

	c.L.Infof("Connecting with %s", c.URL.Scheme+"://*:*@"+c.URL.Host+c.URL.Path+"?"+clean.RawQuery)
	u := connectionString(clean)

	db, err := sql.Open(registeredDriver, u)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open SQL connection")
	}

	c.db = sqlx.NewDb(db, clean.Scheme)
	if err := c.db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "could not ping SQL connection")
	}

	c.L.Infof("Connected to SQL!")

	maxConns := maxParallelism() * 2
	if v := c.URL.Query().Get("max_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.L.Warnf("max_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxConns = int(s)
		}
	}

	maxIdleConns := maxParallelism()
	if v := c.URL.Query().Get("max_idle_conns"); v != "" {
		s, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			c.L.Warnf("max_idle_conns value %s could not be parsed to int: %s", v, err)
		} else {
			maxIdleConns = int(s)
		}
	}

	maxConnLifetime := time.Duration(0)
	if v := c.URL.Query().Get("max_conn_lifetime"); v != "" {
		s, err := time.ParseDuration(v)
		if err != nil {
			c.L.Warnf("max_conn_lifetime value %s could not be parsed to int: %s", v, err)
		} else {
			maxConnLifetime = s
		}
	}

	c.db.SetMaxOpenConns(maxConns)
	c.db.SetMaxIdleConns(maxIdleConns)
	c.db.SetConnMaxLifetime(maxConnLifetime)

	return c.db, nil
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func connectionString(clean *url.URL) string {
	if clean.Scheme == DriverMySQL {
		q := clean.Query()
		q.Set("parseTime", "true")
		clean.RawQuery = q.Encode()
	}

	username := clean.User.Username()
	userinfo := username
	password, hasPassword := clean.User.Password()
	if hasPassword {
		userinfo = url.QueryEscape(userinfo) + ":" + url.QueryEscape(password)
	}
	clean.User = nil
	u := clean.String()
	clean.User = url.UserPassword(username, password)

	if strings.HasPrefix(u, clean.Scheme+"://") {
		u = strings.Replace(u, clean.Scheme+"://", clean.Scheme+"://"+userinfo+"@", 1)
	}
	if clean.Scheme == DriverMySQL {
		u = strings.Replace(u, DriverMySQL+"://", "", -1)

	}
	return u
}

func (c *SQLConnection) registerDriver() (string, error) {
	driverName := c.URL.Scheme
	if c.UseTracedDriver {
		driverName = "instrumented-sql-driver"
		if len(c.options.forcedDriverName) > 0 {
			driverName = c.options.forcedDriverName
		}

		tracingOpts := []instrumentedsql.Opt{instrumentedsql.WithTracer(opentracing.NewTracer(c.AllowRoot))}
		if c.OmitArgs {
			tracingOpts = append(tracingOpts, instrumentedsql.WithOmitArgs())
		}

		switch c.URL.Scheme {
		case DriverPostgreSQL:
			// Why does this have to be a pointer? Because the Open method for postgres has a pointer receiver
			// and does not satisfy the driver.Driver interface.
			sql.Register(driverName,
				instrumentedsql.WrapDriver(&pq.Driver{}, tracingOpts...))
		default:
			return "", fmt.Errorf("unsupported scheme (%s) in DSN", c.URL.Scheme)
		}
	}

	return driverName, nil
}
