#!/usr/bin/python
import time

lcd_columns = 20
lcd_rows    = 4

prevStatus = ""
statusOffset = 0

def readTrainTimes():
	global prevStatus
	global statusOffset
	result = [ " " * lcd_columns ] * lcd_rows
	
	with open("/tmp/l_train_times_and_status_for_lcd.txt") as fh:
		lines = [s.replace("\n", " ") for s in fh.readlines()]

	if len(lines) == 0:
		return result

	result[0:lcd_rows-1] = [x[0:lcd_columns] for x in lines[0:lcd_rows-1]]
	print len(lines)
	status = "" if len(lines) < lcd_rows else lines[lcd_rows-1]	
	prevStatus = status if prevStatus == "" else prevStatus

	if status != "" and status == prevStatus:
		result[lcd_rows-1] = status[0 + statusOffset : lcd_columns + statusOffset]
		statusOffset = statusOffset + 1 if statusOffset < lcd_columns else 0
		if len(result[lcd_rows-1]) < lcd_columns:
			result[lcd_rows-1] = result[lcd_rows-1] + status[0: len(result[lcd_rows-1]) - lcd_columns]
	
	prevStatus = status
	return result


while True:
	print readTrainTimes()
	print "_" * lcd_columns
	time.sleep(0.25)

