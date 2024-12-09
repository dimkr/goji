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

package goji_test

import (
	"fmt"
	"strings"

	"github.com/dimkr/goji"
)

func ExampleJoin() {
	filters := goji.Join(" AND ").
		Add(`sales.price > ?`, 5)

	query, args := goji.Join(" ").
		Add(`SELECT product, SUM(price) FROM sales WHERE`).
		Add(filters).
		Add("GROUP BY product.id HAVING COUNT(DISTINCT stores.id) > ?", 1).
		MustEnd()

	fmt.Printf("%v with %v\n", query, args)

	filters.Add(`sales.price < ?`, 500)

	query, args = goji.Join(" ").
		Add(`SELECT product, SUM(price) FROM sales WHERE`).
		Add(filters).
		Add(`GROUP BY product.id HAVING COUNT(DISTINCT stores.id) > ?`, 1).
		MustEnd()

	fmt.Printf("%v with %v\n", query, args)

	// Output:
	// SELECT product, SUM(price) FROM sales WHERE sales.price > ? GROUP BY product.id HAVING COUNT(DISTINCT stores.id) > ? with [5 1]
	// SELECT product, SUM(price) FROM sales WHERE sales.price > ? AND sales.price < ? GROUP BY product.id HAVING COUNT(DISTINCT stores.id) > ? with [5 500 1]
}

func ExampleBuilder_Add() {
	var a strings.Builder
	a.WriteString("de")
	a.WriteRune('f')

	b := goji.Join(" ").
		Add("x").
		Add("y").
		Add("z")

	query, args := goji.Join("").
		Add("abc", "dog", 123). // string
		Add(&a, "cat", 456).    // fmt.Stringer
		Add(b, "parrot", 789).  // *Builder
		MustEnd()

	fmt.Printf("%v with %v\n", query, args)
	// Output: abcdefx y z with [dog 123 cat 456 parrot 789]
}
