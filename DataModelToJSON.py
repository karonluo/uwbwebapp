
import sys

'''
用于转换本项目的数据库模型为JSON数据模型
Karonsoft DataModel Convert To JSON v0.1
'''


# 目前仅支持 camel 转 pascal


def ProcessFiledName(filed_name, model="pascal", original_model="camel"):
    result = ""
    if model == "pascal":
        if original_model == "camel":
            fileds = filed_name.split("_")
            tmp = ""
            for filed in fileds:
                filed = filed[0].upper() + filed[1:]
                tmp += filed
            result = tmp
    elif model == "camel":
        result = ""
    elif model == "underline":
        result = ""

    return result


def main():
    # if sys.argv[1] != "":
    #     data_model_file_path = sys.argv[1]
    #     if None != sys.argv[2] and sys.argv[2] != "":
    #         rule_model = sys.argv[2]
    #     else:
    #         rule_model = "pascal"
    file_path = sys.argv[1]
    file = open(file_path, "r", encoding="utf-8")
    ctx = file.readlines()
    file.close()
    tmp = "{\r\n"
    for line in ctx:
        line = line.rstrip("\n")
        line = line.rstrip("\r")
        # print('"' + ProcessFiledName(line) + '":' + '"",\r\n')
        tmp += '\t"' + ProcessFiledName(line) + '":' + '"",\r\n'
    tmp = tmp.rstrip(",\r\n")
    tmp += "\r\n}"
    print(tmp)
    file = open(file_path + ".json", "w+")
    file.write(tmp)
    file.close()


if __name__ == "__main__":
    main()
