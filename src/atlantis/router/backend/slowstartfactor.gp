set terminal svg
set output "slowstartfactor.svg"
unset key
set title "slow start factor"
set xlabel "seconds"
set ylabel "factor"
plot "slowstartfactor.csv" using 1:2
