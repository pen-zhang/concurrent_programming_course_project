using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Collections;
using System.Threading;

namespace CombiningTree
{
    class CombiningTree
    {
        Node[] leaf;
        Node[] nodes;

        public CombiningTree(int width)
        {
            nodes = new Node[2 * width - 1];
            nodes[0] = new Node();
            for (int i = 1; i < nodes.Length; i++)
            {
                nodes[i] = new Node(nodes[(i - 1) / 2]);
            }
            leaf = new Node[width];
            for (int i = 0; i < leaf.Length; i++)
            {
                leaf[i] = nodes[nodes.Length - i - 1];
            }
        }

        static long Id()
        {
            return Thread.CurrentThread.ManagedThreadId;
        }

        public int getAndIncrement(int id)
        {
            Stack<Node> stack = new Stack<Node>();
            Node myLeaf = leaf[id % leaf.Length];
            Node node = myLeaf;
            // precombining phase
            try
            {
                while (node.precombine())
                {
                    node = node.parent;
                }
            }
            catch (Exception e) { Console.WriteLine(e); Console.ReadKey(); }
            Node stop = node;
            // combining phase
            int combined = 1;
            for (node = myLeaf; node != stop; node = node.parent)
            {
                try { combined = node.combine(combined); } catch (Exception) { }
                stack.Push(node);
            }


            // operation phase
            int prior = stop.op(combined);
            // distribution phase
            while (stack.Count > 0)
            {
                node = stack.Pop();
                node.distribute(prior);
            }
            return prior;


        }

        public int get()
        {
            return nodes[0].result;
        }
    }
}
