import osmium
import sys
from myanmartools import ZawgyiDetector

detector = ZawgyiDetector()

def check(i,tags):
    if 'name' in tags:
        name = tags['name']
        if ord(name[0]) >= 0x1000 and ord(name[0]) <= 0x109f:
            score = detector.get_zawgyi_probability(name)
            if score > 0.5:
                print(i,name,score)

class NamesHandler(osmium.SimpleHandler):
    def node(self, n):
        check(f"node/{n.id}",n.tags)

    def way(self, w):
        check(f"way/{w.id}",w.tags)

    def relation(self,r):
        check(f"relation/{r.id}",r.tags)

def main(osmfile):
    NamesHandler().apply_file(osmfile)

    return 0

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python %s <osmfile>" % sys.argv[0])
        sys.exit(-1)

    exit(main(sys.argv[1]))