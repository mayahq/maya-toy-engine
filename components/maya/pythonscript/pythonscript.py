import argparse
import zmq
import sys
import json


def main():
    parser = argparse.ArgumentParser(
        description='Cascades component for passing IPs through and logging them to stdout')
    parser.add_argument(
        "-debug", "--debug", help="Enable debug mode", action="store_true")
    parser.add_argument(
        "-json", "--json", help="Print component registration data in JSON", action="store_true")
    parser.add_argument(
        "-port.in", "--port.in", type=str, help="Component input port endpoint")
    parser.add_argument(
        "-port.out", "--port.out", type=str, help="Component out port endpoint")
    args = parser.parse_args()

    if args.json:
        info = dict(
            description="Forwards received IP to the output without any modifications",
            elementary=True,
            inports=(dict(name="IN", type="all", description="Input port for receiving IPs",
                          required=True, addressable=False),
                     dict(name="IN", type="all", description="Input port for receiving IPs",
                          required=True, addressable=False)),
            outports=(dict(name="OUT", type="all",
                           description="Output port for sending IPs",
                           required=True, addressable=False),),
        )
        print(json.dumps(info))
        sys.exit(0)

    data = vars(args)
    in_addr = data.get('port.in')
    out_addr = data.get('port.out')
    if not in_addr or not out_addr:
        parser.print_help()
        sys.exit(1)

    ctx = zmq.Context()

    in_port = ctx.socket(zmq.PULL)
    in_port.bind(in_addr)

    out_port = ctx.socket(zmq.PUSH)
    out_port.connect(out_addr)

    while True:
        if args.debug:
            sys.stdout.write("Waiting for IP...\n")
            sys.stdout.flush()

        parts = in_port.recv_multipart()

        execute_code(parts)

        if args.debug:
            sys.stdout.write("IN Port: Received: %r" % parts)
            sys.stdout.flush()

        out_port.send_multipart(parts)

        if args.debug:
            sys.stdout.write("OUT Port: Forwarded!\n")
            sys.stdout.flush()

def execute_code(parts):
        exec(parts)

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(0)