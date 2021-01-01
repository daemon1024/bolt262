# bolt262

Experimental harness for test262 which can run 100m in 9.58s ⚡😉

## What it is about?

It is cli utility to run [test262 tests](https://github.com/tc39/test262/) for various hosts ( currently tested with node ).
It currently aims to be as fast as possible and optimisations wherever possible and is not comformant with [how tests262 tests should be interpreted](https://github.com/tc39/test262/blob/main/INTERPRETING.md)

## Usage

```sh
bolt262 run [options] <test-file/directory>
```

Run

```sh
bolt262 run help
```

for more info

## Why

- The existing is pretty slow. ( This is almost 100 times faster on rough figures but may get slower once we make it more compliant )
- I wanted to learn more about how test262 works.
- I wanted to dive deep into concurrency and golang in general.

## What's next?

- Make it conformant to how it should be interpreted according to test262.
- Analyse the trace and profiles and analyse the latencies and block inducing code

## Credits

- [ryzokuken](https://github.com/ryzokuken) and [humancalico](https://github.com/humancalico) [sonic262](https://github.com/ryzokuken/sonic262)
- [test262-harness](https://github.com/bterlson/test262-harness)
