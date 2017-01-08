#!/usr/bin/python
import time
import Adafruit_CharLCD as LCD


# Raspberry Pi pin configuration:
lcd_rs        = 27  # Note this might need to be changed to 21 for older revision Pi's.
lcd_en        = 22
lcd_d4        = 25
lcd_d5        = 24
lcd_d6        = 23
lcd_d7        = 18
lcd_backlight = 4

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

# Initialize the LCD using the pins above.
lcd = LCD.Adafruit_CharLCD(lcd_rs, lcd_en, lcd_d4, lcd_d5, lcd_d6, lcd_d7,
                           lcd_columns, lcd_rows, lcd_backlight)

while True:
	data = readTrainTimes()
	x = 0
	y = 0
	for row in data:
		for char in row:
			lcd.set_cursor(y, x)
			lcd.write8(ord(char), True)
			y += 1
		x += 1
	print data
	time.sleep(0.25)

