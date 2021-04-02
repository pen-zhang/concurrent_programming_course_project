using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace BitonicCountingNetwork
{
    class Merger
    {
        private Merger[] _half;

        private Balancer[] _layer;

        private readonly int _width;

        public Merger(int width)
        {
            _width = width;

            _layer = new Balancer[width / 2];
            for (var i = 0; i < width / 2; ++i)
            {
                _layer[i] = new Balancer();
            }

            if (_width > 2)
            {
                _half = new Merger[]
                {
                    new Merger(width / 2),
                    new Merger(width / 2)
                };
            }
        }

        public int Traverse(int input)
        {
            var output = 0;

            if (_width <= 2) return _layer[0].Traverse();

            if (input < _width / 2)
            {
                output = _half[input % 2].Traverse(input / 2);
            }
            else
            {
                output = _half[1 - (input % 2)].Traverse(input / 2);
            }

            return (2 * output) + _layer[output].Traverse();
        }
    }
}
