package mapreduce

import (
	"encoding/json"
	"log"
	"os"
)

// doReduce manages one reduce task: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// You will need to write this function.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTaskNumber) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//

	/**
	 * @author tantexian(https://my.oschina.net/tantexian/blog)
	 * @since 2017/3/3
	 * @params
	 */
	kvMap := make(map[string][]string)
	for i := 0; i < nMap; i++ {
		intermediateFileName := reduceName(jobName, i, reduceTaskNumber)
		file, err := os.Open(intermediateFileName)
		if err != nil {
			log.Fatal("open intermediate file: ", err)
		}
		defer file.Close()

		keyVal := KeyValue{}
		dec := json.NewDecoder(file)
		for {
			err := dec.Decode(&keyVal)
			if err != nil {
				break
			}

			_, ok := kvMap[keyVal.Key]
			if ok {
				// 如果已经存在则append
				kvMap[keyVal.Key] = append(kvMap[keyVal.Key], keyVal.Value)
			} else {
				// 直接添加
				kvMap[keyVal.Key] = []string{keyVal.Value}
			}
		}
	}

	outFileOpen, err := os.Create(outFile)
	if err != nil {
		log.Fatal("open outFile: ", err)
	}
	defer outFileOpen.Close()
	enc := json.NewEncoder(outFileOpen)

	for key, vals := range kvMap {
		// 此处reduceF接受具有相同key的val数组
		enc.Encode(KeyValue{key, reduceF(key, vals)})
	}
}
