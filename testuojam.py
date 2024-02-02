lines = []


def format_time(time):
    return f'{time[0]}{time[1]}:{time[2]}{time[3]}:{time[4]}{time[5]}'



def shrink_line(line):
    line = line.replace("b'", "")
    line = line.split('\\n')[0]
    line = line.replace(line[:6], format_time(line[:6]))
    return line



with open('logas.log', 'r') as f:
    for line in f.readlines():
        shrinked_line = shrink_line(line)
        with open('new_log.log', 'a') as f:
            f.write(f'{shrinked_line} \n')








