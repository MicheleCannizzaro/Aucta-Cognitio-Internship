import requests
import json
import time
import os
import hashlib

#dataloss-prob/component/faults
def data_loss_probability():

    #extract num_osd_in from prometheus
    prometheus_url = 'http://localhost:9090/' #-> Add here the AuctaCognitio Prometheus url  <- LOOK: to change
    osd_in =requests.get(prometheus_url + '/api/v1/query', params={'query': 'dataloss_exporter_pool_dataloss_probability'}) #<- LOOK: to change metric query
    #print(osd_in.text)
    
    osd_in_prometheus = ["osd.1", "osd.5", "osd.4", "osd.15"] #fake_data     <- LOOK:  to comment it out
    osd_in_prometheus = set(osd_in_prometheus)
   
    #extract osd_dump     <- LOOK: to uncomment
    #os.system('ceph osd dump --format=json > osd_dump.json')
    #do an md5 to understand if dump is the same of before or not...
    #is_same= md5_checker("osd_dump.json")

    #if not is_same:
    #    update old file with new one  (probably it is not needed)

    #parse osd_dump to find what osd are in
    with open("../go_server/osd_dump.json") as osd_dump_file: #<- LOOK: to modify path
        osd_dump = json.load(osd_dump_file)

    osd_dump_in = []
       
    for osd in osd_dump["osds"]:
        if osd["in"]==1:
            osd_dump_in.append("osd."+str(osd["osd"]))
    
    #print(osd_dump_in)
    osd_dump_in=set(osd_dump_in)
  
    #check if there is any difference between osd_dump_in and osd_in extracted from prometheus 
    if osd_dump_in != osd_in_prometheus:
        #calculate osd_out                              #LOOK: <- uncomment 
        # if len(osd_in_prometheus)>=len(osd_dump_in):
        #     osd_out= set(osd_in_prometheus)-set(osd_dump_in)
        # else:		
        #     osd_out= set(osd_dump_in)-set(osd_in_prometheus)

        # osd_out=list(osd_out)

        #send query
        osd_out = ["osd.1","osd.2","10.22.22.3", "sv61"] #fake_data  <- LOOK: to comment
        url = 'http://localhost:8081/dataloss-prob/component/faults'

        response = requests.post(url, json= osd_out)
        print(response.text)

#dataloss-prob/component/forecasting
def data_loss_forecasting():
    url = 'http://localhost:8081/dataloss-prob/component/forecasting'
    
    #lack of a way to update osd_lifetime info...

    osd_info_forecasting = [        #fake data  <- LOOK: (that is ok to be fake)
        {
        "osd_name":"osd.1",
        "current_osd_lifetime":60.0,
        "initiation_date":"2019-10-14T02:53:00.000Z"
        },
        {
        "osd_name":"osd.2",
        "current_osd_lifetime":80.0,
        "initiation_date":"2020-10-14T02:53:00.000Z"
        },
        {
        "osd_name":"osd.3",
        "current_osd_lifetime":80.0,
        "initiation_date":"2023-05-01T02:53:00.000Z"
        },
        {
        "osd_name":"osd.2",
        "current_osd_lifetime":79.0,
        "initiation_date":"2023-05-01T02:53:00.000Z"
        }
    ]

    response = requests.post(url, json= osd_info_forecasting)
    print(response.text)

def md5_checker(file_name):

    md5_returned = md5_calculator(file_name)

    if osd_dump_md5 == md5_returned:
        return True
    else:
        return False

def md5_calculator(file_name):
    # Read file and calculate MD5 on its contents 
    with open(file_name, 'rb') as file_to_check:
        # read contents of the file
        data = file_to_check.read()    
        md5_returned = hashlib.md5(data).hexdigest()
    return md5_returned

def init():
    #extract osd_dump     <- LOOK: to uncomment
    #os.system('ceph osd dump --format=json > osd_dump.json') 
    global osd_dump_md5
    osd_dump_md5 = md5_calculator("osd_dump.json")

def stats_update():
    while True:
        data_loss_probability()
        data_loss_forecasting()
        time.sleep(15)

def main():

    global osd_dump_md5

    #init()         #LOOK: <-uncomment
    stats_update()


if __name__ == "__main__":			
	main()