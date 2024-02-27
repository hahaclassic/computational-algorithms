a = float(input("Введите начало диапазона: "))
b = float(input("Введите конец диапазона: "))
step = float(input("Введите шаг: "))
filename = input("Введите название файла: ")

def func(x):
    return x**2

def diff(x):
    return 2*x

def diff2(x):
    return 2

file = open(filename, "w")

while a < b:
    file.write(str(a) + "," + str(func(a)) + "," + str(diff(a)) + "," + str(diff2(a)) + "\n")
    a += step
file.close()