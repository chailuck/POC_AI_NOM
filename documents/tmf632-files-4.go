// internal/middleware/logging.go
package middleware

import (
    "time"
    "github.com/labstack/echo/v4"
    "go.uber.org/zap"
)

func RequestLogging(logger *zap.SugaredLogger) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()

            req := c.Request()
            res := c.Response()

            // Log request details
            logger.Infow("Request started",
                "method", req.Method,
                "uri", req.RequestURI,
                "remote_addr", req.RemoteAddr,
            )

            err := next(c)
            if err != nil {
                c.Error(err)
            }

            // Log response details
            logger.Infow("Request completed",
                "method", req.Method,
                "uri", req.RequestURI,
                "status", res.Status,
                "duration", time.Since(start),
            )

            return err
        }
    }
}

// internal/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
    config := zap.Config{
        Encoding:         "json",
        Level:           zap.NewAtomicLevelAt(zapcore.InfoLevel),
        OutputPaths:     []string{"stdout"},
        ErrorOutputPaths: []string{"stderr"},
        EncoderConfig: zapcore.EncoderConfig{
            MessageKey:   "message",
            LevelKey:     "level",
            TimeKey:      "timestamp",
            EncodeLevel:  zapcore.LowercaseLevelEncoder,
            EncodeTime:   zapcore.ISO8601TimeEncoder,
            EncodeCaller: zapcore.ShortCallerEncoder,
        },
    }

    return config.Build()
}

// internal/validation/validator.go
package validation

import (
    "github.com/go-playground/validator/v10"
    "github.com/your-username/tmf632-service/internal/models"
)

type CustomValidator struct {
    validator *validator.Validate
}

func NewValidator() *CustomValidator {
    return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
    if err := cv.validator.Struct(i); err != nil {
        return err
    }
    return nil
}

func (cv *CustomValidator) ValidateIndividual(individual *models.Individual) error {
    // Add custom validation rules for Individual
    if individual.GivenName == "" {
        return fmt.Errorf("given name is required")
    }
    
    // Validate contact medium
    for _, cm := range individual.ContactMedium {
        if cm.Type == "PhoneContactMedium" && cm.PhoneNumber == "" {
            return fmt.Errorf("phone number is required for PhoneContactMedium")
        }
    }
    
    return nil
}

// deployment/k8s/postgresql.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: tmf632-postgresql
spec:
  serviceName: tmf632-postgresql
  replicas: 1
  selector:
    matchLabels:
      app: tmf632-postgresql
  template:
    metadata:
      labels:
        app: tmf632-postgresql
    spec:
      containers:
      - name: postgresql
        image: postgres:14-alpine
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: tmf632-config
              key: DB_NAME
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: tmf632-secret
              key: DB_USER
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: tmf632-secret
              key: DB_PASSWORD
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        - name: init-script
          mountPath: /docker-entrypoint-initdb.d
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
      - name: init-script
        configMap:
          name: postgres-init-script
---
apiVersion: v1
kind: Service
metadata:
  name: tmf632-postgresql
spec:
  ports:
  - port: 5432
  selector:
    app: tmf632-postgresql
  clusterIP: None
