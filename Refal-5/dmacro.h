
# define chact(arg) { if (arg != actfun) rimp;}
# define gt(m) { if(nst < m) rimp;}
# define lt(m) { if(nst > m) rimp;}
# define eqs(m) { if(nst != m) rimp;}
	/* in EBR sp is decreased to delete the record in st corresponding
		to that break point.	Aug 10 1985. DiMitri Turchin.	*/
# define ebr(m) { sp --; curr_point = m; rd_trace(); nobreak; }
# define nobreak {nel = 3L; b1 = tbel(1); b2 = tbel(2); p = actfun;} 
