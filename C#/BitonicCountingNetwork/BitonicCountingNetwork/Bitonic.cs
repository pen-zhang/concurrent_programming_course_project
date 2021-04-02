using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace BitonicCountingNetwork
{
    public class Bitonic
    {
        private Bitonic[] _half;

        private Merger _merger;

        private readonly int _width;

        public Bitonic(int width)
        {
            _width = width;

            _merger = new Merger(_width);

            if (_width > 2)
            {
                _half = new Bitonic[]
                {
                    new Bitonic(_width / 2),
                    new Bitonic(_width / 2)
                };
            }
        }

        public int Traverse(int input)
        {
            int output = 0;

            if (_width > 2)
            {
                output = _half[input / (_width / 2)].Traverse(input / 2);
            }
            return _merger.Traverse((input >= (_width / 2) ? (_width / 2) : 0) + output);
        }
    }
}
