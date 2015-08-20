## Getting Started
There are a couple of ways to get started with Pact. The easiest method is to [download the compiled binary ]()for your platform and just start using it. 
If however, you're interested in compiling it yourself the following steps should get you there:

1. Clone this repository
1. Run `make`. This should create a `pact` binary which you can use directly like so: `./pact`; you can also copy the binary to `/usr/bin/` or `/usr/local/bin` to use it directly.
1. Run `pact config` to generate a config file. This will also generate a keypair if one does not already exist, and configure the "self" pact which will contain your own public key.
1. You should now be able to run the create and read methods! Try it: `pact create self "some message" | pact read`