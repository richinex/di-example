package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test structures
type TestService interface {
    GetName() string
}

type testServiceImpl struct {
    name string
}

func (t *testServiceImpl) GetName() string {
    return t.name
}

type TestStruct struct {
    Service  TestService `di:"testService"`
    Optional TestService `di:"optionalService"`
    NoTag    TestService
    private  TestService `di:"privateService"`
}

func TestNewContainer(t *testing.T) {
    container := NewContainer()
    assert.NotNil(t, container)
    assert.NotNil(t, container.services)
    assert.NotNil(t, container.log)
}

func TestContainer_Register(t *testing.T) {
    container := NewContainer()

    tests := []struct {
        name      string
        qualifier string
        service   interface{}
        wantErr   bool
    }{
        {
            name:      "valid service",
            qualifier: "testService",
            service:   &testServiceImpl{name: "test"},
            wantErr:   false,
        },
        {
            name:      "nil service",
            qualifier: "nilService",
            service:   nil,
            wantErr:   true,
        },
        {
            name:      "duplicate service",
            qualifier: "testService",
            service:   &testServiceImpl{name: "duplicate"},
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := container.Register(tt.qualifier, tt.service)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                // Verify service was stored
                service, exists := container.services[tt.qualifier]
                assert.True(t, exists)
                assert.Equal(t, tt.service, service)
            }
        })
    }
}

func TestContainer_Resolve(t *testing.T) {
    container := NewContainer()
    testService := &testServiceImpl{name: "test"}

    // Register a test service
    err := container.Register("testService", testService)
    require.NoError(t, err)

    tests := []struct {
        name      string
        qualifier string
        want      interface{}
        wantErr   bool
    }{
        {
            name:      "existing service",
            qualifier: "testService",
            want:      testService,
            wantErr:   false,
        },
        {
            name:      "non-existent service",
            qualifier: "nonExistent",
            want:      nil,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := container.Resolve(tt.qualifier)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}

func TestContainer_InjectStruct(t *testing.T) {
    container := NewContainer()
    testService := &testServiceImpl{name: "test"}

    // Register test service
    err := container.Register("testService", testService)
    require.NoError(t, err)

    tests := []struct {
        name    string
        target  interface{}
        wantErr bool
    }{
        {
            name:    "valid struct pointer",
            target:  &TestStruct{},
            wantErr: false,
        },
        {
            name:    "non-pointer",
            target:  TestStruct{},
            wantErr: true,
        },
        {
            name:    "nil target",
            target:  nil,
            wantErr: true,
        },
        {
            name:    "pointer to non-struct",
            target:  new(string),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := container.InjectStruct(tt.target)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            if ts, ok := tt.target.(*TestStruct); ok {
                // Verify injection
                assert.Equal(t, testService, ts.Service)
                assert.Nil(t, ts.Optional)  // Optional service wasn't registered
                assert.Nil(t, ts.NoTag)     // No tag, shouldn't be injected
                assert.Nil(t, ts.private)   // Private field, can't be injected
            }
        })
    }
}

// TestConcurrency tests thread safety
func TestConcurrency(t *testing.T) {
    container := NewContainer()
    done := make(chan bool)

    // Multiple goroutines registering services
    for i := 0; i < 10; i++ {
        go func(id int) {
            service := &testServiceImpl{name: fmt.Sprintf("service-%d", id)}
            qualifier := fmt.Sprintf("service-%d", id)
            err := container.Register(qualifier, service)
            assert.NoError(t, err)
            done <- true
        }(i)
    }

    // Wait for all goroutines
    for i := 0; i < 10; i++ {
        <-done
    }

    // Verify all services were registered
    assert.Equal(t, 10, len(container.services))
}