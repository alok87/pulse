package controllers

import (
	"net/http"
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"hash/crc32"
	"time"
)

const (
        nexusBasePath = "http://nexus.global.organization.org/nexus/content/repositories/15_OFS/pulse/"
)

func pollNexus(dataPath string, localDataFile string, localNexusFile string) {
for {
		if _, err := os.Stat(localDataFile); err != nil {
			downloadFileFromNexus(localDataFile, dataPath) 
		}		
		
		downloadFileFromNexus(localNexusFile, dataPath) 
	
		if _, err := os.Stat(localDataFile); err != nil {
    			//fmt.Println("File localDataFile not found:", localDataFile)
			continue
		}		

		if _, err := os.Stat(localNexusFile); err != nil {
			//fmt.Println("File localNexusFile not found: ", localNexusFile)
			continue
		}

		h1, err := getHash(localNexusFile)
	 	if err != nil {
	    	fmt.Println("Error getting nexus file's hash", err)
		//continue
		break
	  	}
	  	
	  	h2, err := getHash(localDataFile)
	  	if err != nil {
	    	fmt.Println("Error getting local data file's hash", err)
	    	//continue
		break
	  	}
		
		if h1 != h2 {
			err := os.Remove(localDataFile)
			if err != nil {
				fmt.Println("Error deleting the local data file", err)
				break
			}
			err = os.Rename(localNexusFile, localDataFile)
			if err != nil {
				fmt.Println("Error moving the latest nexus file to local", err)
				break	
			}			
		}
		time.Sleep(100 * time.Millisecond)
		
}
}

func getHash(filename string) (uint32, error) {
  bs, err := ioutil.ReadFile(filename)
  if err != nil {
    return 0, err
  }
  h := crc32.NewIEEE()
  h.Write(bs)
  return h.Sum32(), nil
}


func downloadFileFromNexus(file string, dataPath string) {
for {
	out, err := os.Create(file)
	if err != nil {
		fmt.Println("Error Creating file in local, continuing", err)	
		//panic(err)
	}	
	
	nexusPath := fmt.Sprint(nexusBasePath, dataPath)

	resp, err := http.Get(nexusPath)
	if err != nil {
		fmt.Println("Error Fetching the file from Nexus ", err)
		continue
	}
	
	defer out.Close()	
    	io.Copy(out, resp.Body)
	defer resp.Body.Close()
	break	
}
}
