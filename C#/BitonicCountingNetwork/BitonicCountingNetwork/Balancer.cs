using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;

namespace BitonicCountingNetwork
{
    public class Balancer
    {
        private int _toggle = 1;

        public int Traverse()
        {
            while (true)
            {
                if (1 == Interlocked.Exchange(ref _toggle, 0))
                {
                    return 0;
                }
                if (0 == Interlocked.Exchange(ref _toggle, 1))
                {
                    return 1;
                }
            }
        }
    }
}
