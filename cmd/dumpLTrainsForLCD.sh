#!/bin/sh

outputFilepath="/tmp/l_train_times_and_status_for_lcd.txt"

main() {
	printNextLTrainsForLCD -key "$1" > $outputFilepath && \
	printStatusOfLTrain | grep --after-context 1 'Planned Work' >> $outputFilepath
}

main $*