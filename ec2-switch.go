package main

import (
    "fmt"
    "flag"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
)

func contains(m map[string]string, k string, v string ) bool {
     for key, value := range m{
        if key == k && value == v {
              return true
        }
     }
    return false
}

func act(args string) {
    if (args == "start" || args == "stop" || args == "status"){
        aws_conf := aws.NewConfig()
        sess := session.New(aws_conf)
        svc := ec2.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
        res, err := svc.DescribeInstances(nil)
        if err != nil {
            panic(err)
        }

        for _, r := range res.Reservations {
            for _, i := range r.Instances {
                tags := map[string]string{}
                for _, t := range i.Tags {
                     tags[*t.Key] = *t.Value
                }
                dev := contains(tags, "Env", "dev")
                ec2-switch := contains(tags, "Signal", "true")
                if dev && ec2-switch {
                    if args == "stop" {
                        if (*i.State.Name == "running") {
                            params := &ec2.StopInstancesInput{
                                InstanceIds: []*string{aws.String(*i.InstanceId)},
                            }
                            _, err := svc.StopInstances(params)
                            if err != nil {
                                fmt.Println(err.Error())
                                return
                            }

                            fmt.Println(
                                tags["Name"],
                                *i.InstanceId,
                                "change status: running -> stopped",
                            )
                        }
                    } else if args == "start" {
                        if (*i.State.Name == "stopped") {
                            params := &ec2.StartInstancesInput{
                                InstanceIds: []*string{aws.String(*i.InstanceId)},
                            }
                            _, err := svc.StartInstances(params)
                            if err != nil {
                                fmt.Println(err.Error())
                                return
                            }
                            fmt.Println(
                                tags["Name"],
                                *i.InstanceId,
                                "change status: stopped -> running",
                            )
                        }
                    } else if args == "status" {
                        fmt.Println(
                            tags["Name"],
                            *i.InstanceId,
                            *i.State.Name,
                        )
                    } else {
                      fmt.Println("No instance to change status.")
                    }
                } else {
                  // print instance id staus not changed.
                  // fmt.Println("{ Live Instance: ",*i.InstanceId, "status not changed","}")
                }
            }
        }
    }else{
        fmt.Println("Missing option -o start/stop")
    }
}

func main() {
    var opt = flag.String("o", "blank", "Option to start or stop")
    flag.Parse()

    act(*opt)
}
