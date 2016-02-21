package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/digitalocean/doctl/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/digitalocean/doctl/Godeps/_workspace/src/github.com/digitalocean/godo"

	"github.com/digitalocean/doctl/Godeps/_workspace/src/golang.org/x/oauth2"
)

var SSHCommand = cli.Command{
	Name:    "sshkey",
	Usage:   "SSH Key commands.",
	Aliases: []string{"ssh", "keys"},
	Action:  sshList,
	Subcommands: []cli.Command{
		{
			Name:    "create",
			Usage:   "<name> <path to ssh key(~/.ssh/id_rsa)> Create SSH key.",
			Aliases: []string{"c"},
			Action:  sshCreate,
		},
		{
			Name:    "list",
			Usage:   "List all SSH keys.",
			Aliases: []string{"l"},
			Action:  sshList,
		},
		{
			Name:    "find",
			Usage:   "<name> Find SSH key.",
			Aliases: []string{"f"},
			Action:  sshFind,
		},
		{
			Name:    "destroy",
			Usage:   "[--id | --fingerprint | <name>] Destroy SSH key.",
			Aliases: []string{"d"},
			Action:  sshDestroy,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "id", Usage: "ID for SSH Key. (e.g. 1234567)"},
				cli.StringFlag{Name: "id", Usage: "Fingerprint for SSH Key. (e.g. aa:bb:cc)"},
			},
		},
	},
}

func sshCreate(ctx *cli.Context) {
	if len(ctx.Args()) != 2 {
		fmt.Printf("Must provide name and public key file.\n")
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	file, err := os.Open(ctx.Args()[1])
	if err != nil {
		fmt.Printf("Error opening key file: %s\n", err)
		os.Exit(1)
	}

	keyData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading key file: %s\n", err)
		os.Exit(1)
	}

	createRequest := &godo.KeyCreateRequest{
		Name:      ctx.Args().First(),
		PublicKey: string(keyData),
	}
	key, _, err := client.Keys.Create(createRequest)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	outputKey(*key)
}

func sshList(ctx *cli.Context) {
	if ctx.BoolT("help") == true {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	opt := &godo.ListOptions{}
	keyList := []godo.Key{}

	for {
		keyPage, resp, err := client.Keys.List(opt)
		if err != nil {
			fmt.Printf("Unable to list SSH Keys: %s\n", err)
			os.Exit(1)
		}

		// append the current page's droplets to our list
		for _, d := range keyPage {
			keyList = append(keyList, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			fmt.Printf("Unable to get pagination: %s\n", err)
			os.Exit(1)
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	outputKeys(keyList)
}

func sshFind(ctx *cli.Context) {
	if len(ctx.Args()) != 1 {
		fmt.Printf("Error: Must provide name for Key.\n")
		os.Exit(1)
	}

	name := ctx.Args().First()

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	key, err := FindKeyByName(client, name)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(64)
	}

	outputKey(*key)
}

func sshDestroy(ctx *cli.Context) {
	if ctx.Int("id") == 0 && ctx.String("fingerprint") == "" && len(ctx.Args()) < 1 {
		fmt.Printf("Error: Must provide ID, fingerprint or name for SSH Key to destroy.\n")
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	id := ctx.Int("id")
	fingerprint := ctx.String("fingerprint")
	var key godo.Key
	if id == 0 && fingerprint == "" {
		key, err := FindKeyByName(client, ctx.Args().First())
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(64)
		} else {
			id = key.ID
		}
	} else if id != 0 {
		key, _, err := client.Keys.GetByID(id)
		if err != nil {
			fmt.Printf("Unable to find SSH Key: %s\n", err)
			os.Exit(1)
		} else {
			id = key.ID
		}
	} else {
		key, _, err := client.Keys.GetByFingerprint(fingerprint)
		if err != nil {
			fmt.Printf("Unable to find SSH Key: %s\n", err)
			os.Exit(1)
		} else {
			id = key.ID
		}
	}

	_, err := client.Keys.DeleteByID(id)
	if err != nil {
		fmt.Printf("Unable to destroy SSH Key: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Key %s destroyed.\n", key.Name)
}

type keyOutputable godo.Key

func (a keyOutputable) Headers() []string           { return []string{"ID", "Name", "Fingerprint"} }
func (a keyOutputable) FormatString() string        { return "%d\t%s\t%s\n" }
func (a keyOutputable) RowObject(i int) interface{} { return a }
func (a keyOutputable) RowValues(datum interface{}) []interface{} {
	key := datum.(godo.Key)
	nameLen := len(key.Name)
	if nameLen > 100 {
		nameLen = 100
	}
	return []interface{}{key.ID, key.Name[0:nameLen], key.Fingerprint}
}
func (a keyOutputable) Len() int { return 1 }

func outputKey(key godo.Key) {
	WriteOutputable(keyOutputable(key))
}

type keysOutputable []godo.Key

func (a keysOutputable) Headers() []string           { return []string{"ID", "Name", "Fingerprint"} }
func (a keysOutputable) FormatString() string        { return "%d\t%s\t%s\n" }
func (a keysOutputable) RowObject(i int) interface{} { return a[i] }
func (a keysOutputable) RowValues(datum interface{}) []interface{} {
	key := datum.(godo.Key)
	nameLen := len(key.Name)
	if nameLen > 100 {
		nameLen = 100
	}
	return []interface{}{key.ID, key.Name[0:nameLen], key.Fingerprint}
}
func (a keysOutputable) Len() int { return len(a) }

func outputKeys(keys []godo.Key) {
	WriteOutputable(keysOutputable(keys))
}
