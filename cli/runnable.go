package cli

import (
    "bufio"
    "io"
    "fmt"
    "strconv"
    "errors"
    "github.com/21Bruce/resolved-server/api"
    "os"
)

var (
    ErrNoName = errors.New("name required. use -n flag")
    ErrMulArg = errors.New("too many arguments for a flag")
    ErrNoArg = errors.New("too few arguments for a flag")
)

type APIFlag struct {
    API      []api.API
    Name string
} 

type ResolvedCLI struct {
    API         api.API 
    In          io.Reader
    Out         io.Writer
    Err         io.Writer
    parseCtx    ParseCtx
}

func validateSearch(in map[string][]string) (string, int, error){
    var err error = nil     

    if in["n"] == nil {
        err = ErrNoName
    }
    if len(in["n"]) > 1 {
        err = ErrMulArg
    }
    if len(in["n"]) == 0 {
        err = ErrNoArg
    }
    if (in["l"] != nil) && (len(in["l"]) > 1) {
        err = ErrMulArg
    }
    if (in["l"] != nil) && (len(in["l"]) == 0) {
        err = ErrNoArg
    }

    if err != nil {
        return "", 0, err
    }

    name := in["n"][0]
    limit := 0
    if in["l"] != nil {
        limitRes, err := strconv.Atoi(in["l"][0])
        if err != nil {
            return "", 0 , err
        }
        limit = limitRes
    }
    return name, limit, nil

}

func (c *ResolvedCLI) handleSearch(in map[string][]string) (string, error) {
    name, limit, err := validateSearch(in)
    if err != nil {
        return "", err
    }
    searchParams := api.SearchParam{Name: name, Limit: limit}
    resp, err := c.API.Search(searchParams)
    if err != nil {
        return "", err
    }
    return resp.ToString(), nil
}

func (c *ResolvedCLI) handleQuit(in map[string][]string) (string, error) {
    os.Exit(0)
    return "", nil
}

func (c *ResolvedCLI) handleHelp(in map[string][]string) (string, error) {
    helpStr := "Commands: \n"
    for _, cmd := range c.parseCtx.Commands {
        helpStr += "\t" + cmd.Name 
        for _, flag := range cmd.Flags {
            helpStr += " [-" + flag.Name + "]"
        }
        helpStr += ": "+ cmd.Description + "\n"
        for _, flag := range cmd.Flags {
            helpStr += "\t\t[-" + flag.Name + "]: "  + flag.Description + "\n"
        }
    }

    return helpStr, nil 
}

//func validateRats(in map[string][]string) (string, int, error){
//
//}

func (c *ResolvedCLI) handleRats(in map[string][]string) (string, error) {
    return "", nil
}

func (c *ResolvedCLI) initParseCtx() {
    searchCommand := Command{
        Name: "search",
        Description: "Finds restaurant info",
        Flags: []Flag{
            Flag{
                Name: "n",
                Description: "This flag is required. It takes one text input, the name of the restaurant",
            },
            Flag{
                Name: "l",
                Description: "This flag is optional. It takes one number input, the max amount of results to return",
            },
        },
        Handler: c.handleSearch,
    }

    ratsCommand := Command{
        Name: "rats",
        Description: "Reserve At Time Scheduler. Sends a reservation request at the specified time",
        Flags: []Flag{
            Flag{
                Name: "v",
                Description: "This flag is required. It takes one number input, the id of the restaurant(use search command)",
            },
            Flag{
                Name: "l",
                Description: "This flag is required. It takes a list of military times in hh:mm:ss format to try to reserve at",
            },
            Flag{
                Name: "d",
                Description: "This flag is required. It takes a day in yyyy:mm:dd format to try to reserve at",
            },
            Flag{
                Name: "r",
                Description: "This flag is required. It takes a day and time in yyyy:mm:dd:hh:mm:ss format to begin sending reservation requests at",
            },
            Flag{
                Name: "s",
                Description: "This flag is optional. It takes a text name as input and uses the venueID from the first result of that search",
            },
        },
        Handler: c.handleRats,
    }

    quitCommand := Command{
        Name: "quit",
        Description: "Exits the CLI",
        Flags: []Flag{},
        Handler: c.handleQuit,
    }

    exitCommand := Command{
        Name: "exit",
        Description: "Exits the CLI",
        Flags: []Flag{},
        Handler: c.handleQuit,
    }

    helpCommand := Command{
        Name: "help",
        Description: "Displays helpful info about commands",
        Flags: []Flag{},
        Handler: c.handleHelp,
    }

    c.parseCtx = ParseCtx{
        OpenDelim: "[",
        CloseDelim: "]",
        Commands: []Command{
            searchCommand,
            ratsCommand,
            quitCommand,
            exitCommand,
            helpCommand,
        },
    }
}

func (c *ResolvedCLI) Run() {
    c.initParseCtx()
    scanner := bufio.NewScanner(c.In)
    fmt.Fprintln(c.Out, "Welcome to the Resolved CLI! For Help type 'help'") 
    for {
        fmt.Fprint(c.Out, "resolved(0.1.0)>> ") 
        scanner.Scan()
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(c.Err, err);
        }
        result, err := c.parseCtx.Parse(scanner.Text()) 
        if err != nil {
            fmt.Fprint(c.Err, "ERROR: ")
            fmt.Fprintln(c.Err, err)
        } else  {
            fmt.Fprintln(c.Out, result) 
        }
    }
}


