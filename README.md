# observer

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/imkira/go-observer/blob/master/LICENSE.txt)
[![GoDoc](https://godoc.org/github.com/imkira/go-observer?status.svg)](https://godoc.org/github.com/imkira/go-observer)
[![Build Status](http://img.shields.io/travis/imkira/go-observer.svg?style=flat)](https://travis-ci.org/imkira/go-observer)
[![Coverage](https://codecov.io/gh/imkira/go-observer/branch/master/graph/badge.svg)](https://codecov.io/gh/imkira/go-observer)
[![codebeat badge](https://codebeat.co/badges/28bdd579-8b34-4940-a3e0-35ac52794a42)](https://codebeat.co/projects/github-com-imkira-go-observer)
[![goreportcard](https://goreportcard.com/badge/github.com/imkira/go-observer)](https://goreportcard.com/report/github.com/imkira/go-observer)

observer is a [Go](http://golang.org) package that aims to simplify the problem
of channel-based broadcasting of events from one or more publishers to one or
more observers.

# Problem

The typical quick-and-dirty approach to notifying a set of observers in go is
to use channels and call each in a for loop, like the following:

```go
for _, channel := range channels {
  channel <- value
}
```

There are two problems with this approach:

- The broadcaster blocks every time some channel is not ready to be written to.
- If the broadcaster blocks for some channel, the remaining channels will not
  be written to (and therefore not receive the event) until the blocking
  channel is finally ready.
- It is O(N). The more observers you have, the worse this loop will behave.

Of course, this could be solved by creating one goroutine for each channel so
the broadcaster doesn't block. Unfortunately, this is heavy and
resource-consuming. This is especially bad if you have events being raised
frequently and a considerable number of observers.

# Approach

The way observer package tackles this problem is very simple. For every event,
a state object containing information about the event, and a channel is
created. State objects are managed using a singly linked list structure: every
state points to the next. When a new event is raised, a new state object is
appended to the list and the channel of the previous state is closed (this
helps notify all observers that the previous state is outdated).

Package observer defines 2 concepts:

- Property: An object that is continuously updated by one or more publishers.
- Stream: The list of values a property is updated to. For every property
update, that value is appended to the list in the order they happen, and is
only discarded when you advance to the next value.

# Memory Usage

The amount of memory used for one property is not dependent on the number of
observers. It should be proportional to the number of value updates since the
value last obtained by the slowest observer. As long as you keep advancing all
your observers, garbage collection will take place and keep memory usage
stable.

# How to Use

First, you need to install the package:

```
go get -u github.com/imkira/go-observer
```

Then, you need to include it in your source:

```go
import "github.com/imkira/go-observer"
```

The package will be imported with ```observer``` as name.

The following example creates one property that is updated every second by one
or more publishers, and observed by one or more observers.

## Documentation

For advanced usage, make sure to check the
[available documentation here](http://godoc.org/github.com/imkira/go-observer).

## Example: Creating a Property

The following code creates a property with initial value ```1```.

```go
val := 1
prop := observer.NewProperty(val)
```

After creating the property, you can pass it around to publishers or
observers as you want.

## Example: Publisher

The following code represents a publisher that increments the value of the
property by one every second.

```go
val := 1
for {
  time.Sleep(time.Second)
  val += 1
  fmt.Printf("will publish value: %d\n", val)
  prop.Update(val)
}
```

Note:

- Property is goroutine safe: you can use it concurrently from multiple
goroutines.

## Example: Observer

The following code represents an observer that prints the initial value of a
property and waits indefinitely for changes to its value. When there is a
change, the stream is advanced and the current value of the property is
printed.

```go
stream := prop.Observe()
val := stream.Value().(int)
fmt.Printf("initial value: %d\n", val)
for {
  select {
    // wait for changes
    case <-stream.Changes():
      // advance to next value
      stream.Next()
      // new value
      val = stream.Value().(int)
      fmt.Printf("got new value: %d\n", val)
  }
}
```

Note:

- Stream is not goroutine safe: You must create one stream by calling
  ```Property.Observe()``` or ```Stream.Clone()``` if you want to have
  concurrent observers for the same property or stream.

## Example

Please check
[examples/multiple.go](https://github.com/imkira/go-observer/blob/master/examples/multiple.go)
for a simple example on how to use multiple observers with a single updater.
