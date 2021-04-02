using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;
using System.Diagnostics;

namespace CombiningTree
{
    class Test
    {
        static CombiningTree tree;
        static int TH = 12, NUM = 100;
        static int width;
        // tree: combining tree where threads do increment ops
        // TH: number of threads
        // NUM: number of increment ops each thread performs

        static void Main(String[] args)
        {
            width = (int) Math.Ceiling((decimal)TH / 2);
            tree = new CombiningTree(width);
            var stopWatch = new Stopwatch();


            Log(width + "-leaves Combining tree.");
            Log("Starting " + TH + " threads doing increments ...");

            stopWatch.Start();

            Parallel.For(0, TH, (index) =>
            {   
                DateTime start = DateTime.Now;
                for (int i = 0; i < NUM; i++)
                   tree.getAndIncrement(index);
                DateTime stop = DateTime.Now;
                Log(index + ": done in " + (stop - start).TotalSeconds + "s");
            });

            stopWatch.Stop();
            var ts = stopWatch.Elapsed;
            Log("Total: " + tree.get());
            Log("Total time: " + String.Format("{0:00}.{1:00}", ts.Seconds, ts.Milliseconds / 10)+'s');

            Console.ReadKey();
        }

        static void Log(String x)
        {
            Console.WriteLine(x);
        }

}
}

