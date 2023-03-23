#!/bin/python3

import json
def set_interface(pathsNode, pathNode, urlMethod="get"):
	interface_dict = None
	try:
		interface_dict = {}
		interface_dict["url"] = pathNode
		interface_dict["method"] = urlMethod	
		interface_dict["display_name"] = pathsNode[pathNode][urlMethod]["summary"]
	except:
		interface_dict = None
		pass
	return interface_dict
	
def main():
	file = open("./uwbfuncpages.openapi.json", "rb")
	interface_dict_list=[]
	obj = json.load(file)
	file.close()
	paths = obj["paths"]
	for path in paths:
		interface_dict = set_interface(paths, path, "get")
		if interface_dict is not None:
			interface_dict_list.append(interface_dict)
		interface_dict = set_interface(paths, path, "post")
		if interface_dict is not None:
			interface_dict_list.append(interface_dict)
		interface_dict = set_interface(paths, path, "delete")
		if interface_dict is not None:
			interface_dict_list.append(interface_dict)
		interface_dict = set_interface(paths, path, "put")
		if interface_dict is not None:
			interface_dict_list.append(interface_dict)
	filew = open("./sysfuncpages.json", "wb")
	filew.write( bytes(json.dumps(interface_dict_list, ensure_ascii=False), encoding="utf-8") )
	filew.close()

if __name__ == "__main__":
	main()
