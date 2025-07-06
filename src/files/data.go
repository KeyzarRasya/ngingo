package files

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/KeyzarRasya/ngingo/src/balancer"
)

type DataWriteRead interface {
	Write(records [][]string)	error
	Read()						([]balancer.VarStat, error)
}

type DataCPU struct {
	filepath 	string;
	VarStat 	balancer.VarStat
}

func NewDataCPU(filepath string, varstat balancer.VarStat) DataCPU {
	return DataCPU{filepath: filepath, VarStat:  varstat}
}

func (d *DataCPU) Write(records [][]string) error {
	f, err := os.OpenFile(d.filepath, os.O_APPEND | os.O_WRONLY, 0644)

	if err != nil {
		return err;
	}



	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err;
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func (d *DataCPU) Read() ([]balancer.VarStat, error) {
	var varStats []balancer.VarStat;

	f, err := os.OpenFile(d.filepath, os.O_RDONLY, 0644);


	if err != nil {
		return nil, err;
	}

	r := csv.NewReader(f)

	records, err := r.ReadAll()

	if err != nil {
		return nil, err;
	}

	for _, record := range records {
		varstat := d.VarStat.Clone()

		port, err := strconv.ParseUint(record[0], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid port at line X: %w", err)
		}
		usage, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float at line X: %w", err)
		}
		varstat.SetPortVarStat(uint16(port), record[1], usage)
		varStats = append(varStats, varstat)

	}
	fmt.Println(varStats)

	return varStats, nil
}
