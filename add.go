package main

import (
	"fmt"
	"log"
	"reflect"
	"github.com/urfave/cli"
	"github.com/satori/go.uuid"
)

type AddJob struct {
	BaseJob
	uuid bool
}

var addFlags = []cli.Flag {
	cli.BoolFlag {
		Name: "uuid",
		Usage: "use UUID",
	},
}

func Add(c *cli.Context) error {
	runBenchmark(c, reflect.TypeOf(AddJob{}))
    return nil
}

func (job *AddJob) Prep(c *cli.Context) bool {
	if job.GetVerbose() >= 1 {
		log.Printf("worker[%d]: prepare\n", job.wid)
	}
	job.uuid = c.Bool("uuid")
	err := job.ldap.Bind(c.String("D"), c.String("w"))
	if err != nil {
		log.Fatal("bind error: ", err)
		return false
	}
	return true
}

func (job *AddJob) Request() bool {
	var cn string
	id, err := uuid.NewV1()
	if job.uuid {
		cn = id.String()
	} else {
		cn = fmt.Sprintf("%d-%d", job.wid, job.count)
	}
	dn := fmt.Sprintf("cn=%s,%s", cn, job.baseDN)
	attrs := map[string][]string{
		"objectClass": {"person"},
		"cn": {cn},
		"sn": {"sn"},
		"userPassword": {"secret"},
	}
	err = job.ldap.Add(dn, attrs)
	if err != nil {
		log.Printf("add error: %s", err)
		return false
	}
	return true
}
