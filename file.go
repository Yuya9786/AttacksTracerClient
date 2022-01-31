package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	pb "github.com/Yuya9786/AttacksTracerClient/protobuf"
	"github.com/pkg/errors"
)

type FileEntity struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Annotations FileAnnotations `json:"annotations"`
}

type FileAnnotations struct {
	BootID     int    `json:"boot_id"`
	Date       string `json:"cf:date"`
	Version    int    `json:"version"`
	Epoch      int    `json:"cf:epoch"`
	Taint      string `json:"cf:taint"`
	Pathname   string `json:"pathname"`
	ObjectID   string `json:"object_id"`
	Jiffies    string `json:"cf:jiffies"`
	ObjectType string `json:"object_type"`
	MachineID  string `json:"cf:machine_id"`
}

var target = "/home/fedora/secret.txt"

func read(db *sql.DB, c pb.MalwareSimulatorClient, relation *Relation) error {
	if relation.Annotations.ToType != "file" || relation.Annotations.FromType != "task" {
		return errors.New("format is not matched")
	}

	name, err := getFileName(db, relation.To)
	if err != nil {
		return errors.Wrap(err, "failed to fileCheck")
	}

	if name == nil {
		// path is not found
		println(relation.To, "path not found")
		return nil
	}

	if *name != target {
		// this is not target
		return nil
	}

	app, err := task(db, c, relation.From)
	if err != nil {
		return err
	}

	_, err = readFile(c, app, *name)
	if err != nil {
		return err
	}

	return nil
}

func getFileName(db *sql.DB, file string) (*string, error) {
	query := fmt.Sprintf("select * from data where record->>'from'='%s';", file)
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "faile to query")
	}
	defer rows.Close()

	var path *string
	for rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan rows")
		}

		var relation Relation
		if err = json.Unmarshal(data.Record, &relation); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}

		if relation.Annotations.ToType == "path" {
			path = &relation.To
			break
		}

		if relation.Annotations.RelationType == "version_entity" {
			return getFileName(db, relation.To)
		}
	}

	if path == nil {
		fmt.Println("path is not found", file)
		return nil, nil
	}

	query = fmt.Sprintf("select * from data where record->>'id'='%s';", *path)
	rows, err = db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "faile to query")
	}
	defer rows.Close()

	if rows.Next() {
		var data Data
		err = rows.Scan(&data.Tag, &data.Time, &data.Record)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan rows")
		}

		var entity FileEntity
		if err = json.Unmarshal(data.Record, &entity); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}

		pathname := entity.Annotations.Pathname

		return &pathname, nil
	}

	return nil, errors.New("pathname not found")
}
