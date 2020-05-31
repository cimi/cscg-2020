from argparse import ArgumentParser
from validate import validate
from solve import solve

parser = ArgumentParser()
parser.add_argument("-c", "--command", dest="command",
                    help="one of: process, train, validate, solve", metavar="FILE")
def main():
  args = parser.parse_args()
  if args.command == 'validate':
    validate("model-homogenous2.hdf5")
    validate("model-homogenous-500.hdf5")
  elif args.command == 'solve':
    solve()
  else: 
    print(args)

if __name__ == "__main__":
    main()
