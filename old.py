from time import sleep
from datetime import datetime
import socket
import sys

now = datetime.now()
dt_string = now.strftime("%Y%m%d_%H%M%S") + ".log"


class AnalyzerReader:
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    
    def connect(self):
        try:
            server_address = (self.host, self.port)
            self.sock.connect(server_address)
        except:
            print("Cannot connect to Cube, check if software is running")

    def disconnect(self):
        self.sock.close()

    
    def send_message(self, text):
        msg = f'{text}\r\n'
        encoded_msg = bytes(msg, 'ascii')
        self.sock.sendall(encoded_msg)

    
    def get_status(self):
        self.send_message('?STS')
        sleep(0.1)
        while True:
            data = self.sock.recv(4096)
            if not data:
                break
            print(f'Received: {data}')
            return

    def get_status_continuous(self):
        while True:
            data = self.sock.recv(4096)
            if not data:
                break
            print(f'Received: {data}')
            with open(dt_string, 'a') as f:
                #f.write('\n' + str(data))
                #now.strftime("%H%M%S\t")
                now = datetime.now()
                #f.write(now.strftime("%H%M%S\t") + str(data.decode()))
                f.write('\n' + now.strftime("%H%M%S\t") + str(data))
            return

    def get_name(self, position):
        self.send_message(f'?NAM {position}')
        sleep(0.1)
        while True:
            data = self.sock.recv(4096)
            if not data:
                break
            print(f'Received: {data}')
            return
        
    def strt(self):
        self.send_message('STRT')
        sleep(0.1)


    def seqon(self):
        self.send_message('SEQON')
        sleep(0.1)

    def sleepoff(self):
        self.send_message('SLEEPOFF')
        sleep(0.1)

    def sleepon(self):
        self.send_message('SLEEPON')
        sleep(0.1)

    def OpenLog(self):
        with open(dt_string, 'w') as f:
            f.write(dt_string + '\n')

                
if __name__ == '__main__':
    ar = AnalyzerReader('localhost', 1984)
    ar.OpenLog()
    ar.connect()
    #ar.sleepoff()
    #ar.sleepon()
    #sleep(1)
    #ar.get_status()
    #ar.get_name(1)
    #ar.strt()
    #ar.seqon()
    while 1:
        ar.get_status_continuous()
        #print('in a loop!!')
        sleep(0.01)
    ar.disconnect()
