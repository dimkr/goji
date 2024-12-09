/*
 * This file is part of goji.
 *
 * Copyright (c) 2024 Dima Krasner
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package goji implements a string joiner with a fluent interface.
package goji

import (
	"fmt"
	"strings"
)

// Builder is an alternative to [strings.Builder] that supports chaining and parameterized expressions.
//
// It can be used as a poor man's SQL query builder, for use with [database/sql].
type Builder struct {
	inner strings.Builder
	delim string
	args  []any
	err   error
}

// Join returns a new [Builder] which joins parameterized expressions with a given delimiter.
func Join(delim string) *Builder {
	return &Builder{delim: delim}
}

// End returns the built string and array of parameters.
func (b *Builder) End() (string, []any, error) {
	return b.inner.String(), b.args, b.err
}

// MustEnd is like [Builder.End] but panics on error.
func (b *Builder) MustEnd() (string, []any) {
	if b.err != nil {
		panic("goji: " + b.err.Error())
	}

	return b.inner.String(), b.args
}

func (b *Builder) setErr(err error) {
	if b.err != nil {
		b.err = err
	}
}

func (b *Builder) Write(p []byte) (int, error) {
	if b.inner.Len() > 0 {
		p = append([]byte(b.delim), p...)
	}

	n, err := b.inner.Write(p)
	if err != nil {
		b.setErr(err)
	}

	return n, err
}

// Add appends a parameterized expression.
//
// exp can be anything, including another [*Builder] or a [fmt.Stringer].
func (b *Builder) Add(exp any, arg ...any) *Builder {
	var err error
	switch v := exp.(type) {
	case string:
		if b.inner.Len() > 0 {
			v = b.delim + v
		}

		if _, err = b.inner.WriteString(v); err != nil {
			b.setErr(err)
		} else {
			b.args = append(b.args, arg...)
		}

	case *Builder:
		s, more, err := v.End()
		if err != nil {
			b.setErr(err)
		} else {
			if b.inner.Len() > 0 {
				s = b.delim + s
			}

			if _, err = b.inner.WriteString(s); err != nil {
				b.setErr(err)
			} else {
				b.args = append(b.args, more...)
				b.args = append(b.args, arg...)
			}
		}

	case fmt.Stringer:
		s := v.String()
		if b.inner.Len() > 0 {
			s = b.delim + s
		}

		if _, err = b.inner.WriteString(s); err != nil {
			b.setErr(err)
		} else {
			b.args = append(b.args, arg...)
		}

	default:
		var err error
		if b.inner.Len() > 0 {
			_, err = fmt.Fprintf(&b.inner, "%s%v", b.delim, v)
		} else {
			_, err = fmt.Fprintf(&b.inner, "%v", v)
		}

		if err != nil {
			b.setErr(err)
		} else {
			b.args = append(b.args, arg...)
		}
	}

	return b
}
