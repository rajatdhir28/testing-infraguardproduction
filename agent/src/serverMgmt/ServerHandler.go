

package serverMgmt
// version No 1 dated :- 03-Apr-2017
import (
    "strings"
    "encoding/json"
    "net/http"
    "agentUtil"
    "stringUtil"
    "fileUtil"
  _ "fmt" // for unused variable issue
    "io/ioutil"
    "fmt"
    "net/url"
    
    //"fmt"    
)

/*{"fieldCount":0,"affectedRows":1,"insertId":44,
  "serverStatus":2,"warningCount":0,"message":"","protocol41":true,"changedRows":0}*/


//const baseUrl = "https://ojf489mkrc.execute-api.us-west-2.amazonaws.com/dev/registerserver"
//const baseUrl = "https://09q09swczl.execute-api.ap-southeast-1.amazonaws.com/prod/registerserver" 

func DoServerRegnProcess(urlForServerRegn string) (string){
 
    if(len(urlForServerRegn) == 0){
      fileUtil.WriteIntoLogFile("Missing url for server regn.")
      return "1"
    }

    url := urlForServerRegn + getQueryString()
    fmt.Println("\n\nServer full url for regn = : ",url) 
    fileUtil.WriteIntoLogFile("ServerHandler.DoServerRegnProcess(). Going to hit url = : "+url)


    res, err := http.Get(url)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN 27. Msg = : "+err.Error())
        return "1"
    }
    body, err := ioutil.ReadAll(res.Body)
    
    if err != nil {
      fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN. 32. Msg = : "+err.Error())
      return "1"
    }
    var data interface{} 
    err = json.Unmarshal(body, &data)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN 38. Msg = : "+err.Error())
        return "1"
    }

   infraGuardResp, _ := data.(map[string]interface{})
    
    var affectedRows float64
    affectedRows = -1
          
    if val, ok := infraGuardResp["affectedRows"].(float64); ok {
      affectedRows = val
    } else {
      errorMsg := "ServerHandler.DoServerRegnProcess() LN 59. Unable to cast into float64"
      fileUtil.WriteIntoLogFile(errorMsg)

    }
   if(affectedRows > 0){
      return "0"
    }else{
     return "1"
    }

}


func getQueryString()(string){
   serverIp := agentUtil.ExecComand("curl -s http://169.254.169.254/latest/meta-data/local-ipv4", "ServerHandler.go 74")
   hostName := agentUtil.ExecComand("hostname", "ServerHandler.go 75")
   globalHostName := agentUtil.ExecComand("curl -s http://169.254.169.254/latest/meta-data/public-hostname", "ServerHandler.go 500")
   instanceId := agentUtil.ExecComand("curl -s http://169.254.169.254/latest/meta-data/instance-id", "ServerHandler.go 501")
   region := agentUtil.ExecComand("curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone | sed 's/.$//'", "ServerHandler.go 502")
   mac := agentUtil.ExecComand("curl -s http://169.254.169.254/latest/meta-data/mac", "ServerHandler.go 503")
   vpcId := agentUtil.ExecComand(" curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/"+mac+"/vpc-id", "ServerHandler.go 504")
	
   serverIp = strings.TrimSpace(serverIp)
   hostName = strings.TrimSpace(hostName)
   globalHostName = strings.TrimSpace(globalHostName)
   
   instanceId = strings.TrimSpace(instanceId)
   region = strings.TrimSpace(region)
   vpcId = strings.TrimSpace(vpcId)
   
   userList := agentUtil.ExecComand("cat /etc/passwd | grep '/home' | cut -d: -f1", "ServerHandler.go 71")
   userList2 := strings.Split(userList,"\n")

  
   max := 500
   if(len(userList2) < max){
    max = len(userList2)
   }

  users := ""
  for i := 0; i <  max; i++ {
    if(len(users) ==0){
      users = userList2[i]
    }else{
      users = users +","+userList2[i]
    }
 
  }
 users = strings.TrimSpace(users)
 
 // Remove last Comma
 if(strings.Contains(users, ",")){
     users = string(users[0:(len(users)-1)])
  }

 /*
    Read Command line arguments given at the time of agentInstaller.sh execution
 */
 //var sName, pId, licenseKey string
 var sName, pId, licenseKey string
 
 sName = "sName"
 pId = "5"
 licenseKey = "lKey"

 if(fileUtil.IsFileExisted("/tmp/serverInfo.txt")){
  args := stringUtil.SplitData(fileUtil.ReadFile("/tmp/serverInfo.txt", false), ":")
  if(len(args) == 3){
    sName = url.QueryEscape(args[0]);
    pId = args[1];
    licenseKey = args[2];
    agentUtil.ExecComand("rm -r /tmp/serverInfo.txt", "ServerHandler.go 106")
  }
 }

 
 cpuDetails := agentUtil.ExecComand("lscpu", "ServerHandler.go 105")
 cpuDetails = stringUtil.FindKey(cpuDetails)

 kernelDetails := agentUtil.ExecComand("cat /etc/*-release", "ServerHandler.go 114")
 kernelDetails = stringUtil.FindKey(kernelDetails)

//url.QueryEscape("how can I do this")
 qryStr := "?serverName="+sName+"&serverIp="+serverIp+"&hostName="+hostName+"&instanceId="+instanceId+"&region="+region+"&vpcId="+vpcId+"&globalHostName="+globalHostName+"&projectId="+pId+"&userList="+users+"&licenseKey="+licenseKey
 qryStr = strings.Replace(qryStr, "\n","",-1)
 return qryStr
}


