package httpexpect

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain_Basic(t *testing.T) {
	t.Run("clone", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.clone()

		assert.NotSame(t, chain1, chain2)
		assert.NotSame(t, chain1.context.Path, chain2.context.Path)

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())
	})

	t.Run("enter_leave", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.enter("test")

		assert.NotSame(t, chain1, chain2)
		assert.NotSame(t, chain1.context.Path, chain2.context.Path)

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())

		chain2.leave()
	})

	t.Run("enter_leave_fail", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.enter("test")

		chain2.fail(mockFailure())

		assert.False(t, chain1.failed())
		assert.True(t, chain2.failed())

		assert.False(t, chain1.treeFailed())
		assert.True(t, chain2.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, flagFailed)

		chain2.leave() // propagates failure to parents

		assert.True(t, chain1.failed())
		assert.True(t, chain2.failed())

		assert.True(t, chain1.treeFailed())
		assert.True(t, chain2.treeFailed())

		chain1.assertFlags(t, flagFailed)
		chain2.assertFlags(t, flagFailed)
	})

	t.Run("clone_enter_leave_fail", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.clone()
		chain3 := chain2.clone()
		chain3e := chain3.enter("test")

		chain3e.fail(mockFailure())

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.False(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())
		assert.False(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, 0)
		chain3.assertFlags(t, 0)
		chain3e.assertFlags(t, flagFailed)

		chain3e.leave() // propagates failure to parents

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.True(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.True(t, chain1.treeFailed())
		assert.True(t, chain2.treeFailed())
		assert.True(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, flagFailedChildren)
		chain2.assertFlags(t, flagFailedChildren)
		chain3.assertFlags(t, flagFailed)
		chain3e.assertFlags(t, flagFailed)
	})

	t.Run("two_branches", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.clone()
		chain2e := chain2.enter("test")
		chain3 := chain2.clone()
		chain3e := chain3.enter("test")

		chain2e.fail(mockFailure())
		chain3e.fail(mockFailure())

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.True(t, chain2e.failed())
		assert.False(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())
		assert.True(t, chain2e.treeFailed())
		assert.False(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, 0)
		chain2e.assertFlags(t, flagFailed)
		chain3.assertFlags(t, 0)
		chain3e.assertFlags(t, flagFailed)

		chain2e.leave() // propagates failure to parents
		chain3e.leave() // propagates failure to parents

		assert.False(t, chain1.failed())
		assert.True(t, chain2.failed())
		assert.True(t, chain2e.failed())
		assert.True(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.True(t, chain1.treeFailed())
		assert.True(t, chain2.treeFailed())
		assert.True(t, chain2e.treeFailed())
		assert.True(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, flagFailedChildren)
		chain2.assertFlags(t, flagFailed|flagFailedChildren)
		chain2e.assertFlags(t, flagFailed)
		chain3.assertFlags(t, flagFailed)
		chain3e.assertFlags(t, flagFailed)
	})

	t.Run("set_root_1", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.clone()
		chain2.setRoot()
		chain3 := chain2.clone()
		chain3e := chain3.enter("test")

		chain3e.fail(mockFailure())

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.False(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())
		assert.False(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, 0)
		chain3.assertFlags(t, 0)
		chain3e.assertFlags(t, flagFailed)

		chain3e.leave() // propagates failure to parents

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.True(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.False(t, chain1.treeFailed())
		assert.True(t, chain2.treeFailed())
		assert.True(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, flagFailedChildren)
		chain3.assertFlags(t, flagFailed)
		chain3e.assertFlags(t, flagFailed)
	})

	t.Run("set_root_2", func(t *testing.T) {
		chain1 := newMockChain(t)
		chain2 := chain1.clone()
		chain3 := chain2.clone()
		chain3.setRoot()
		chain3e := chain3.enter("test")

		chain3e.fail(mockFailure())

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.False(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())
		assert.False(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, 0)
		chain3.assertFlags(t, 0)
		chain3e.assertFlags(t, flagFailed)

		chain3e.leave() // propagates failure to parents

		assert.False(t, chain1.failed())
		assert.False(t, chain2.failed())
		assert.True(t, chain3.failed())
		assert.True(t, chain3e.failed())

		assert.False(t, chain1.treeFailed())
		assert.False(t, chain2.treeFailed())
		assert.True(t, chain3.treeFailed())
		assert.True(t, chain3e.treeFailed())

		chain1.assertFlags(t, 0)
		chain2.assertFlags(t, 0)
		chain3.assertFlags(t, flagFailed)
		chain3e.assertFlags(t, flagFailed)
	})
}

func TestChain_Panics(t *testing.T) {
	t.Run("nil_reporter", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = newChainWithDefaults("test", nil)
		})
	})

	t.Run("set_request_twice", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.setRequest(&Request{})

		assert.Panics(t, func() {
			opChain.setRequest(&Request{})
		})
	})

	t.Run("set_response_twice", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.setResponse(&Response{})

		assert.Panics(t, func() {
			opChain.setResponse(&Response{})
		})
	})

	t.Run("leave_without_enter", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		assert.Panics(t, func() {
			chain.leave()
		})
	})

	t.Run("leave_on_parent", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		_ = chain.enter("foo")

		assert.Panics(t, func() {
			chain.leave()
		})
	})

	t.Run("double_leave", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.leave()

		assert.Panics(t, func() {
			opChain.leave()
		})
	})

	t.Run("enter_after_leave", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.leave()

		assert.Panics(t, func() {
			opChain.enter("bar")
		})
	})

	t.Run("replace_without_enter", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		assert.Panics(t, func() {
			chain.replace("foo")
		})
	})

	t.Run("replace_after_leave", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.leave()

		assert.Panics(t, func() {
			opChain.replace("bar")
		})
	})

	t.Run("replace_empty_path", func(t *testing.T) {
		chain := newChainWithDefaults("", newMockReporter(t))

		opChain := chain.enter("")

		assert.Panics(t, func() {
			opChain.replace("bar")
		})
	})

	t.Run("replace_empty_aliased_path", func(t *testing.T) {
		chain := newChainWithDefaults("", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.setAlias("")

		assert.Panics(t, func() {
			opChain.replace("bar")
		})
	})

	t.Run("fail_without_enter", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		assert.Panics(t, func() {
			chain.fail(mockFailure())
		})
	})

	t.Run("fail_after_leave", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.leave()

		assert.Panics(t, func() {
			opChain.fail(mockFailure())
		})
	})

	t.Run("clone_after_leave", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.leave()

		assert.Panics(t, func() {
			_ = opChain.clone()
		})
	})

	t.Run("alias_after_leave", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")
		opChain.leave()

		assert.Panics(t, func() {
			opChain.setAlias("bar")
		})
	})

	t.Run("setters_after_leave", func(t *testing.T) {
		setterFuncs := []func(chain *chain){
			func(chain *chain) {
				chain.setRoot()
			},
			func(chain *chain) {
				chain.setSeverity(SeverityLog)
			},
			func(chain *chain) {
				chain.setRequestName("")
			},
			func(chain *chain) {
				chain.setRequest(nil)
			},
			func(chain *chain) {
				chain.setResponse(nil)
			},
		}

		for _, setter := range setterFuncs {
			chain := newChainWithDefaults("test", newMockReporter(t))

			opChain := chain.enter("foo")
			opChain.leave()

			assert.Panics(t, func() {
				setter(opChain)
			})
		}
	})

	t.Run("invalid_assertion", func(t *testing.T) {
		chain := newChainWithDefaults("test", newMockReporter(t))

		opChain := chain.enter("foo")

		assert.Panics(t, func() {
			opChain.fail(AssertionFailure{
				Type: AssertionType(9999),
			})
		})
	})
}

func TestChain_Env(t *testing.T) {
	t.Run("newChainWithConfig_non_nil_env", func(t *testing.T) {
		env := NewEnvironment(newMockReporter(t))

		chain := newChainWithConfig("root", Config{
			AssertionHandler: &mockAssertionHandler{},
			Environment:      env,
		}.withDefaults())

		assert.Same(t, env, chain.env())
	})

	t.Run("newChainWithConfig_nil_env", func(t *testing.T) {
		chain := newChainWithConfig("root", Config{
			AssertionHandler: &mockAssertionHandler{},
			Environment:      nil,
		}.withDefaults())

		assert.NotNil(t, chain.env())
	})

	t.Run("newChainWithDefaults", func(t *testing.T) {
		chain := newChainWithDefaults("root", newMockReporter(t))

		assert.NotNil(t, chain.env())
	})
}

func TestChain_Root(t *testing.T) {
	t.Run("newChainWithConfig_non_empty_path", func(t *testing.T) {
		chain := newChainWithConfig("root", Config{
			AssertionHandler: &mockAssertionHandler{},
		}.withDefaults())

		assert.Equal(t, []string{"root"}, chain.context.Path)
	})

	t.Run("newChainWithConfig_empty_path", func(t *testing.T) {
		chain := newChainWithConfig("", Config{
			AssertionHandler: &mockAssertionHandler{},
		}.withDefaults())

		assert.Equal(t, []string{}, chain.context.Path)
	})

	t.Run("newChainWithDefaults_non_empty_path", func(t *testing.T) {
		chain := newChainWithDefaults("root", newMockReporter(t))

		assert.Equal(t, []string{"root"}, chain.context.Path)
	})

	t.Run("newChainWithDefaults_empty_path", func(t *testing.T) {
		chain := newChainWithDefaults("", newMockReporter(t))

		assert.Equal(t, []string{}, chain.context.Path)
	})
}

func TestChain_Path(t *testing.T) {
	path := func(c *chain) string {
		return strings.Join(c.context.Path, ".")
	}

	rootChain := newChainWithDefaults("root", newMockReporter(t))

	assert.Equal(t, "root", path(rootChain))

	opChain1 := rootChain.enter("foo")

	assert.Equal(t, "root", path(rootChain))
	assert.Equal(t, "root.foo", path(opChain1))

	opChain2 := opChain1.enter("bar")

	assert.Equal(t, "root", path(rootChain))
	assert.Equal(t, "root.foo", path(opChain1))
	assert.Equal(t, "root.foo.bar", path(opChain2))

	opChain2Clone := opChain2.clone()
	opChain3 := opChain2Clone.enter("baz")

	assert.Equal(t, "root", path(rootChain))
	assert.Equal(t, "root.foo", path(opChain1))
	assert.Equal(t, "root.foo.bar", path(opChain2))
	assert.Equal(t, "root.foo.bar", path(opChain2Clone))
	assert.Equal(t, "root.foo.bar.baz", path(opChain3))

	opChain1r := opChain1.replace("xxx")
	opChain3r := opChain3.replace("yyy")

	assert.Equal(t, "root", path(rootChain))
	assert.Equal(t, "root.foo", path(opChain1))
	assert.Equal(t, "root.foo.bar", path(opChain2))
	assert.Equal(t, "root.foo.bar", path(opChain2Clone))
	assert.Equal(t, "root.foo.bar.baz", path(opChain3))
	assert.Equal(t, "root.xxx", path(opChain1r))
	assert.Equal(t, "root.foo.bar.yyy", path(opChain3r))
}

func TestChain_AliasedPath(t *testing.T) {
	path := func(c *chain) string {
		return strings.Join(c.context.Path, ".")
	}
	aliasedPath := func(c *chain) string {
		return strings.Join(c.context.AliasedPath, ".")
	}

	t.Run("enter_and_leave", func(t *testing.T) {
		rootChain := newChainWithDefaults("root", newMockReporter(t))

		assert.Equal(t, "root", path(rootChain))
		assert.Equal(t, "root", aliasedPath(rootChain))

		c1 := rootChain.enter("foo")
		assert.Equal(t, "root.foo", path(c1))
		assert.Equal(t, "root.foo", aliasedPath(c1))

		c2 := c1.enter("bar")
		assert.Equal(t, "root.foo.bar", path(c2))
		assert.Equal(t, "root.foo.bar", aliasedPath(c2))

		c2.setAlias("alias1")
		assert.Equal(t, "root.foo.bar", path(c2))
		assert.Equal(t, "alias1", aliasedPath(c2))

		c3 := c2.enter("baz")
		assert.Equal(t, "root.foo.bar.baz", path(c3))
		assert.Equal(t, "alias1.baz", aliasedPath(c3))

		c3.leave()
		assert.Equal(t, "root.foo.bar.baz", path(c3))
		assert.Equal(t, "alias1.baz", aliasedPath(c3))
	})

	t.Run("set_empty", func(t *testing.T) {
		rootChain := newChainWithDefaults("root", newMockReporter(t))

		assert.Equal(t, "root", path(rootChain))
		assert.Equal(t, "root", aliasedPath(rootChain))

		rootChain.setAlias("")
		assert.Equal(t, "root", path(rootChain))
		assert.Equal(t, "", aliasedPath(rootChain))
	})
}

func TestChain_Handler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		handler := &mockAssertionHandler{}

		chain := newChainWithConfig("test", Config{
			AssertionHandler: handler,
		}.withDefaults())

		opChain := chain.enter("test")
		opChain.leave()

		assert.NotNil(t, handler.ctx)
		assert.Nil(t, handler.failure)
	})

	t.Run("failure", func(t *testing.T) {
		handler := &mockAssertionHandler{}

		chain := newChainWithConfig("test", Config{
			AssertionHandler: handler,
		}.withDefaults())

		opChain := chain.enter("test")
		opChain.fail(mockFailure())
		opChain.leave()

		assert.NotNil(t, handler.ctx)
		assert.NotNil(t, handler.failure)
	})
}

func TestChain_Severity(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		handler := &mockAssertionHandler{}

		chain := newChainWithConfig("test", Config{
			AssertionHandler: handler,
		}.withDefaults())

		opChain := chain.enter("test")
		opChain.fail(mockFailure())
		opChain.leave()

		assert.NotNil(t, handler.failure)
		assert.Equal(t, SeverityError, handler.failure.Severity)
	})

	t.Run("error", func(t *testing.T) {
		handler := &mockAssertionHandler{}

		chain := newChainWithConfig("test", Config{
			AssertionHandler: handler,
		}.withDefaults())

		chain.setSeverity(SeverityError)

		opChain := chain.enter("test")
		opChain.fail(mockFailure())
		opChain.leave()

		assert.NotNil(t, handler.failure)
		assert.Equal(t, SeverityError, handler.failure.Severity)
	})

	t.Run("log", func(t *testing.T) {
		handler := &mockAssertionHandler{}

		chain := newChainWithConfig("test", Config{
			AssertionHandler: handler,
		}.withDefaults())

		chain.setSeverity(SeverityLog)

		opChain := chain.enter("test")
		opChain.fail(mockFailure())
		opChain.leave()

		assert.NotNil(t, handler.failure)
		assert.Equal(t, SeverityLog, handler.failure.Severity)
	})
}
