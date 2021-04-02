# Chapter12 Counting

This project focuses on the foundation of distributed coordination. 

According to the book and online materials, I finish implementation of the combining tree, bitonic counting network algorithm in Java and C# for now. There are still problems in  Go version. 

We will talk about the node structure and four operation phase of combining tree to make a counter visited by multiple threads. Based on the idea, the author introduced the counting network replacing the single counter with multiple counters in the combining tree. 

Once the counting network is built, we can get a sorting network if we change the balancer in counting network to a comparator. But the sorting network typically work will for small in-memory data sets. 

