[Benchmarks](./benchmark_test.go)

```
Benchmark_Iterator/Constant-12         	594181771	         1.92 ns/op	       0 B/op	       0 allocs/op
Benchmark_Iterator/Linear-12           	684941280	         1.70 ns/op	       0 B/op	       0 allocs/op
Benchmark_Iterator/LinearRate-12       	677859252	         1.72 ns/op	       0 B/op	       0 allocs/op
Benchmark_Iterator/Exponential-12      	735283754	         1.57 ns/op	       0 B/op	       0 allocs/op
Benchmark_Iterator/ExponentialRate-12  	425531461	         2.81 ns/op	       0 B/op	       0 allocs/op
```

```
Benchmark_IteratorWithMaxRetries/Constant-12            	306885199	         3.83 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithMaxRetries/Linear-12              	319124492	         3.77 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithMaxRetries/LinearRate-12          	324360441	         3.67 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithMaxRetries/Exponential-12         	321683217	         3.63 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithMaxRetries/ExponentialRate-12     	310087774	         3.77 ns/op	       0 B/op	       0 allocs/op
```

```
Benchmark_IteratorWithJitter/Constant-12         	46091090	        26.0 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithJitter/Linear-12           	47933851	        24.8 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithJitter/LinearRate-12       	46190442	        24.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithJitter/Exponential-12      	47949747	        24.9 ns/op	       0 B/op	       0 allocs/op
Benchmark_IteratorWithJitter/ExponentialRate-12  	44407602	        24.6 ns/op	       0 B/op	       0 allocs/op
```

```
Benchmark_Iterable/Constant-12         	79993066	        14.1 ns/op	       8 B/op	       1 allocs/op
Benchmark_Iterable/Linear-12           	71026089	        17.5 ns/op	      16 B/op	       1 allocs/op
Benchmark_Iterable/LinearRate-12       	70618558	        18.1 ns/op	      16 B/op	       1 allocs/op
Benchmark_Iterable/Exponential-12      	92273622	        14.0 ns/op	       8 B/op	       1 allocs/op
Benchmark_Iterable/ExponentialRate-12  	67006164	        18.2 ns/op	      16 B/op	       1 allocs/op
```

```
Benchmark_IterableWithMaxRetries/LinearRate-12         	23091488	        48.4 ns/op	      48 B/op	       2 allocs/op
Benchmark_IterableWithMaxRetries/Exponential-12        	27702172	        44.4 ns/op	      40 B/op	       2 allocs/op
Benchmark_IterableWithMaxRetries/ExponentialRate-12    	25327890	        48.3 ns/op	      48 B/op	       2 allocs/op
Benchmark_IterableWithMaxRetries/Constant-12           	28570951	        44.1 ns/op	      40 B/op	       2 allocs/op
Benchmark_IterableWithMaxRetries/Linear-12             	25531969	        47.9 ns/op	      48 B/op	       2 allocs/op
```

```
Benchmark_IterableWithJitter/Constant-12         	24490144	        48.4 ns/op	      40 B/op	       2 allocs/op
Benchmark_IterableWithJitter/Linear-12           	23751087	        51.3 ns/op	      48 B/op	       2 allocs/op
Benchmark_IterableWithJitter/LinearRate-12       	23751886	        52.3 ns/op	      48 B/op	       2 allocs/op
Benchmark_IterableWithJitter/Exponential-12      	24998749	        49.0 ns/op	      40 B/op	       2 allocs/op
Benchmark_IterableWithJitter/ExponentialRate-12  	23377867	        52.5 ns/op	      48 B/op	       2 allocs/op
```

```
Benchmark_Trier/Constant-12         	     100	  10424801 ns/op	     218 B/op	       4 allocs/op
Benchmark_Trier/Linear-12           	     100	  10323035 ns/op	     225 B/op	       4 allocs/op
Benchmark_Trier/LinearRate-12       	     100	  10294533 ns/op	     225 B/op	       4 allocs/op
Benchmark_Trier/Exponential-12      	     100	  10311722 ns/op	     216 B/op	       4 allocs/op
Benchmark_Trier/ExponentialRate-12  	     100	  10269813 ns/op	     225 B/op	       4 allocs/op
```

```
Benchmark_TrierWithMaxRetries/Constant-12         	     100	  10301934 ns/op	     250 B/op	       5 allocs/op
Benchmark_TrierWithMaxRetries/Linear-12           	     100	  10297511 ns/op	     256 B/op	       5 allocs/op
Benchmark_TrierWithMaxRetries/LinearRate-12       	     100	  10356039 ns/op	     257 B/op	       5 allocs/op
Benchmark_TrierWithMaxRetries/Exponential-12      	     100	  10271901 ns/op	     250 B/op	       5 allocs/op
Benchmark_TrierWithMaxRetries/ExponentialRate-12  	     100	  10300833 ns/op	     257 B/op	       5 allocs/op
```

```
Benchmark_TrierWithJitter/Constant-12            	      97	  10561697 ns/op	     248 B/op	       5 allocs/op
Benchmark_TrierWithJitter/Linear-12              	     100	  10651914 ns/op	     258 B/op	       5 allocs/op
Benchmark_TrierWithJitter/LinearRate-12          	     100	  10670898 ns/op	     256 B/op	       5 allocs/op
Benchmark_TrierWithJitter/Exponential-12         	      99	  10718809 ns/op	     248 B/op	       5 allocs/op
Benchmark_TrierWithJitter/ExponentialRate-12     	     100	  10536359 ns/op	     256 B/op	       5 allocs/op
```

```
Benchmark_TrierWithMaxRetriesJitter/Constant-12         	     100	  10726373 ns/op	     282 B/op	       6 allocs/op
Benchmark_TrierWithMaxRetriesJitter/Linear-12           	     100	  10704179 ns/op	     288 B/op	       6 allocs/op
Benchmark_TrierWithMaxRetriesJitter/LinearRate-12       	     100	  10475046 ns/op	     288 B/op	       6 allocs/op
Benchmark_TrierWithMaxRetriesJitter/Exponential-12      	     100	  10405025 ns/op	     280 B/op	       6 allocs/op
Benchmark_TrierWithMaxRetriesJitter/ExponentialRate-12  	     100	  10548585 ns/op	     288 B/op	       6 allocs/op
```
