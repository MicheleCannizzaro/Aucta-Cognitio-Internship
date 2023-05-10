import requests
import json
import time
import os
import hashlib
import sys
import logging
import re

#dataloss-prob/component/faults
def pool_data_loss_probability():

    try:
        #---extract num_osd_in from prometheus---   
        if os.path.exists("prometheus_url.txt"):
            with open('prometheus_url.txt') as f:
                prometheus_url = f.readlines()[0].split('\n')[0]
        else:
            logging.warning("prometheus_url.txt not found")
            

        osd_in =requests.get(prometheus_url + '/api/v1/query', params={'query': 'ceph_osd_in'}).json()  #returns a dictionary
        osd_in = json.dumps(osd_in) #return a well formed json
        
        #parsing dell'output di requests in json 
        osds_in_prometheus = []

        if osd_in["status"] == "success":

            for result in osd_in["data"]["result"]:

                if result["value"][1] == "1":
                    osds_in_prometheus.append(result["metric"]["ceph_daemon"])

        else:
            logging.critical("Reading prometheus osd_in reported failure")

        #osds_in_prometheus = ["osd.1", "osd.2"] #fake_data
        osds_in_prometheus = set(osds_in_prometheus)
   
        #---extracting osd_dump from cluster--- 
        #to understand if osd_dump is the same of before or not... and writing it only if there are differences
        output_md5_dump=os.popen('sudo ceph osd dump --format=json | md5sum').read()

        pattern = '([a-z0-9]*)'
        osd_dump_md5 = re.match(pattern, output_md5_dump).group()
        
        if not md5_checker(osd_dump_md5):  #is osd_dump the same? if not...
            #need to update osd_dump.json file
                #extracting osd_dump and overwriting the existing one
            os.system('sudo ceph osd dump --format=json > osd_dump.json')
            print("osd_dump.json updated\n")

            global osd_dump_md5_from_file 
            osd_dump_md5_from_file = md5_calculator("osd_dump.json")
            
            print ("(new) osd_dump_md5_from_file: "+osd_dump_md5_from_file)
        
        if os.path.exists("osd_dump.json"):
            
            #parse osd_dump.json to find what osd are in
            with open("osd_dump.json") as osd_dump_file:
                osd_dump = json.load(osd_dump_file)

            osds_dump_in = []
            
            for osd in osd_dump["osds"]:
                if osd["in"]==1:
                    osds_dump_in.append("osd."+str(osd["osd"]))
            
            #print(osds_dump_in)
            osds_dump_in=set(osds_dump_in)
        
            #--- check if there is any difference between osds_dump_in and osd_in extracted from prometheus---
            
            if osds_dump_in != osds_in_prometheus:
                #calculate osds_out                              
                if len(osds_dump_in)>len(osds_in_prometheus):
                    osds_out= set(osds_dump_in)-set(osds_in_prometheus) #some osd is down
                    osds_out=list(osds_out)

                    print(f"osds_out: {osds_out}")

                    #send query
                    #osds_out = ["osd.1"] #fake_data  <- LOOK: to comment
                    url = 'http://localhost:8081/dataloss-prob/component/faults'

                    response = requests.post(url, json= osds_out)
                    print(response.text)

                else:		
                    #osds_out= set(osds_in_prometheus)-set(osds_dump_in)
                    print("Prometheus shows more osds with status 'in' than dump")
                    osds_out=set()                       
                    osds_out=list(osds_out)
        
            print(f"osds_dump_in-> {len(osds_dump_in)} osds_in_prometheus-> {len(osds_in_prometheus)} OUT->{len(osds_out)}")
            print("------------------------------------------")
        
        else:
            logging.warning("osd_dump.json not found")

    except (requests.exceptions.JSONDecodeError,UnboundLocalError, json.decoder.JSONDecodeError, TypeError) as ex:
         logging.critical("Error in making request to Prometheus metric endpoint "+str(ex)+"\n")        

#dataloss-prob/component/forecasting
def pool_data_loss_forecasting():
    url = 'http://localhost:8081/dataloss-prob/component/forecasting'
    
    #lack of a way to update osd_lifetime info...
    with open ('osds_infos_fake_data.json') as f:
        osd_info_forecasting = json.load(f)

    response = requests.post(url, json= osd_info_forecasting)
    print(response.text)
    print("------------------------------------------")

def md5_checker(osd_dump_md5):

    global osd_dump_md5_from_file 
    osd_dump_md5_from_file = md5_calculator("osd_dump.json")
    
    print("osd_dump_md5_from_file: "+ osd_dump_md5_from_file)
    print("osd_dump_md5: "+ osd_dump_md5)

    if osd_dump_md5_from_file == osd_dump_md5:
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
    #extract osd_dump
    os.system('sudo ceph osd dump --format=json > osd_dump.json') 
    
    #refers to the global variable
    global osd_dump_md5_from_file
    
    try:
        osd_dump_md5_from_file = md5_calculator("osd_dump.json")
    except FileNotFoundError:
        logging.critical("osd_dump.json was not created")

def stats_update():
    while True:
       
        print("---pool data loss probability---")
        pool_data_loss_probability()
       
        print("---pool data loss forecasting---")
        pool_data_loss_forecasting()
        time.sleep(15)

def main():
    print("------------------------------------------")
    global osd_dump_md5_from_file

    init()       
    stats_update()


if __name__ == "__main__":
    try:			
        main()
    except KeyboardInterrupt:
        print("\nTermination...\nBye")
        try:
            sys.exit(130)
        except SystemExit:
            os._exit(130)
    except (requests.exceptions.ConnectionError):
        print("Error: Unable to contact the Server")