# Go sample project

Sample repo to quickstart golang applications. Plz note that this is designed to demonstrate common used features, not best practices. E.g. the custom flagset usage here is mostly superficial, in a program like this it's better to just use default flagset. Same applies to packages - a program as small as this typically can all be done in main package, no need to split anything.

## Functionality

Sample functionality includes:

- Parsing data from stdin or file
- Serializing data to json or xml and posting it to pastebin (or stdout)
- Running a server with a simple `mood` configuration

## Environment setup

1. Install [gvm](https://github.com/moovweb/gvm)
1. For first-time setup use binary install. Note that `--default` is necessary so that the version is available in PATH:

    ```zsh
    gvm install go1.20.3 --prefer-binary --default
    ```
1. Create new project or use existing (compatible) project (one you can run with your default golang version)

    ```sh
    mkdir -p sample-go
    cd peruna
    go mod init github.com/megakoresh/sample-go
    ```

1. After doing all of this, install the go extension for vscode and configure the following in the vscode settings (json):

    ```json
    {
        //...
        "go.goroot": "~/.gvm/gos/go1.20.3",
    }
    ```

    Unfortunately this is needed as there is no way to currently marry VSCode with gvm. If you want to switch your default go version, you will need to update this setting also. With IntelliJ the IDE usually does all of these things for you (trade-off is that you don't know how it does that, so if something goes wrong it's hard to debug). 

## Resources

- <https://gobyexample.com/> - the best place to look for common use-case implementation
- <https://pkg.go.dev/> - standard library docs

## Tips for newcomers

- Code is split to logical packages declared at top of file
    - For monorepo with multiple go apps use go [workspaces](https://go.dev/doc/tutorial/workspaces)
- Only single main package and one `main()` function possible per module
- Variables and functions that start with upper case are available in other packages
- To control how things are parsed between data formats use tags in structs

    ```go
    type DataItem struct {
        Format string `json:"-" xml:"-"` // omit

        Peruna   string `json:"peruna" xml:"peruna"`
        Porkkana string `json:"porkkana" xml:"porkkana"`
    }
    ```
- Golang is not object oriented, but it does have types and methods. Latter are basically syntactic sugar - it's like passing the object as first argument to method, it makes your code easier to work with. Also the combnation of functions that you implement "on" your custom data type is an interface and you can implement some common interfaces just by creating a function for your datatype and it will immediately work with most other golang contexts.

    ```go
    func (di DataItem) String() string {
        return fmt.Sprintf("%s %s", di.Peruna, di.Porkkana)
    }
    ```
- Variable assignment and creation are explicit in golang. There is shorthand creation + assignment `:=` and explicit creation with `var` keyword. Use `var` if you must declare variable before assignment. Otherwise use `:=`.

    ```go
    package peruna

    import (
        "github.com/sok/skaupat"
        "github.com/sok/kasvikset"
    )

    var (
        noutoPiste skaupat.NoutoPiste // creates empty skaupat.NoutoPiste to be shared between all functions in the package
    )

    func calculatePrice() int64 {
        price := 0
        for _, tuote := range noutoPiste.Tuoteet {
            price += tuote.Amount * tuote.PricePerUnit
        }
        return price
    }

    func main() {
        peruna := kasvikset.NewPeruna()
        noutoPiste.LoadVegetable(peruna)
    }
    ```

- Serializing and deserializing in golang implements a well-defined and ubiquitous model for transferring data between different formats, which is strongly typed at all times and easy to use.

    ```go
    // MarshalJSON implements the Marschaller interface so now if anyone calls `json.Marshall` on an instance of NoutoPiste, it will call your method
    func (np skaupat.NoutoPiste) MarshalJSON() ([]byte, error) {
        logger.Fatalf("not implemented") // not implementing for sake of clarity
    }

    func main() {
        fileContentsFull, err := os.ReadFile("/home/stan/peruna.json")
        if err != nil {
            logger.Fatalf("%v", err)
        }
        var peruna kasvikset.Peruna
        // if you already have full data content in memory, use Marshaller
        err := json.Unmarshal(fileContentsFull, &peruna)
        if err != nil {
            logger.Fatalf("Could not parse peruna.json because of %v", err)
        }
        var perunat []kasvikset.Peruna
        // possibly big file, so lets not load it all to memory, ok?
        fd, err := os.Open("/home/stan/way_too_many_perunat.json")
        // defer puts a function call to end of function like a stack. Makes it less likely you will forget to close open handlers
        defer fd.Close()
        if err != nil {
            logger.Fatalf("%v", err)
        }
        dec := json.NewDecoder(fd)
        // if you are reading from a data stream (e.g. network request or file descriptor) use decoder
		if err := dec.Decode(&perunat); err != nil {
			return nil, err
		}
        noutoPiste.LoadVegetable(peruna)
    }
    ```

- If you need to relinquish control of the thead (e.g. waiting for a network request), use goroutines and channels. Remember that you can specify channel direction when declaring a function to tell the compiler if you are expecting to read or write to it.

    ```go
    // we specify chane<- so the compiler knows getData must only write to the channel, not read. So compiler will throw error now if we try to read from it and we'll have to fix it instead of pulling hair out debugging at runtime
    func getData(c chan<- []byte) {
        // implement your IO-bound logic here
    }

    func main() {
        c := make(chan []byte)
        go getData(c)
        logger.Fprintln("This line will execute immedately")
        d := <-c // this line will block until data is available in channel
        logger.Fprintf("%s", d) // now data is avalable
    }
    ```

    Goroutines can be called safely from loops, unlike normal threads they are not necessarily going to be scheduled at the OS level. Still, excersie caution when doing this.
- Golang has pretty much every common use-case covered by the standard library - from mocking requests to websockets. Before writing some implementation, plz consult the best practices
- If you have multiple channels you are expecting some data in, golang has special construct called `select` to monitor channels. I found that it's usually best to have all the `selects` on the main goroutine - that way you get better control of program output.
    - Same thing goes for outputting from other goroutines: IMO it's best to avoid that and just return something to channel whenever a spawned goroutine needs to output something.

    ```go
    select {
        case result := <-myChannel:
        logger.Printf("Received something from my channel: %s", result)
        case err := <-errChannel:
        logger.Fatalf("Received error on error channel: %s", err)
    }
    ```
