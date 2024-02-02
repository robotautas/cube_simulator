from time import sleep
import socket


class AnalyzerReader:
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.connect()

    def connect(self):
        try:
            server_address = (self.host, self.port)
            self.sock.connect(server_address)
        except:
            print("Cannot connect to Cube, check if software is running")

    def disconnect(self):
        self.sock.close()

    def msg(self, text: str) -> str | None:
        msg = f"{text}\r\n"
        encoded_msg = bytes(msg, "ascii")
        self.sock.sendall(encoded_msg)
        print("sent message: ", encoded_msg)
        sleep(0.1)
        # if msg startswith '?', then response is expected, otherwise it's a command
        if text[0] == "?":
            while True:
                data = self.sock.recv(128)

                decoded_data = data.decode()

                # if not data or '\n' in decoded_data:
                return decoded_data
        else:
            return

    def get_status_continuous(self):
        while True:
            data = self.sock.recv(
                4096,
            )
            if not data:
                break
            print(f"Received: {data}")
            return

    def get_name(self, position):
        return self.msg(f"?NAM {position}")

    def get_weight(self, position):
        return self.msg(f"?WGH {position}")

    def get_sinx(self):
        return self.msg(f"?SINX")

    def get_c(self, position):
        return self.msg(f"?PCT {position} C")

    def get_n(self, position):
        return self.msg(f"?PCT {position} N")

    def strt(self):
        self.msg("STRT")
        sleep(0.1)

    def seqon(self):
        self.msg("SEQON")
        sleep(0.1)


if __name__ == "__main__":
    ar = AnalyzerReader("localhost", 1984)
    status = ar.msg("?STS")
    print("sent message 1")
    print(status)

    sleep(1)

    status = ar.msg("?Hello")
    print("sent message 2")
    print(status)

    status = ar.msg("?STS")
    print("sent message 2")
    print(status)

    ar.msg("START")
    while 1:
        ar.get_status_continuous()
        sleep(1)

    # ar.disconnect()
