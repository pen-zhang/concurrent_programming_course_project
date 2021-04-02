using System;
using System.Threading.Tasks;
using System.Diagnostics;
using System.Threading;

namespace BitonicCountingNetwork
{
    class Test
    {
        static void Main(string[] args)
        {
            const int width = 4;
            var bitonic = new Bitonic(width);
            var counters = new int[width];

            const int tokenCount = 1000;
            var tokens = new int[tokenCount];

            var rand = new Random();
            var randLock = new object();
            var stopWatch = new Stopwatch();

            Parallel.For(0, tokenCount, (i) =>
            {
                int next;
                lock (randLock)
                {
                    next = rand.Next(width);
                }
                tokens[i] = next;
            });

            stopWatch.Start();
            Parallel.For(0, tokenCount, (i) =>
            {
                Interlocked.Increment(ref counters[bitonic.Traverse(tokens[i])]);
            });
            stopWatch.Stop();

            var traversing = stopWatch.Elapsed.TotalSeconds;

            for (var i = 0; i < width; ++i)
            {
                Console.WriteLine($"Output: {i} Count: {counters[i]}");
            }

            Console.WriteLine($"Time to traverse the network: {traversing.ToString()} s");

            Console.ReadKey();
        }
    }
}
