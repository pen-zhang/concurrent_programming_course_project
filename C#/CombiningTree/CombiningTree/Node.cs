using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Diagnostics;
using System.Threading;

namespace CombiningTree
{
	class Node
	{
		enum CStatus { IDLE, FIRST, SECOND, RESULT, ROOT };
		bool locked;
		CStatus cStatus;
		int firstValue, secondValue;
		public int result;
		public Node parent;

		public Node()
		{
			cStatus = CStatus.ROOT;
			locked = false;
		}
		public Node(Node myParent)
		{
			parent = myParent;
			cStatus = CStatus.IDLE;
			locked = false;
		}

		public bool precombine()
		{
			lock (this)
			{
				while (locked) Monitor.Wait(this);

				switch (cStatus)
				{
					case CStatus.IDLE:
						cStatus = CStatus.FIRST;
						return true;
					case CStatus.FIRST:
						locked = true;
						cStatus = CStatus.SECOND;
						return false;
					case CStatus.ROOT:
						return false;
					default:
						throw new Exception("unexpected Node state" + cStatus);
				}
			}
		}

		public int combine(int combined)
		{
			lock (this)
			{
				while (locked) Monitor.Wait(this);

				locked = true;
				firstValue = combined;
				switch (cStatus)
				{
					case CStatus.FIRST:
						return firstValue;
					case CStatus.SECOND:
						return firstValue + secondValue;
					default:
						throw new Exception("unexpected Node state " + cStatus);
				}
			}
		}

		public int op(int combined)
		{
			lock (this)
			{
				switch (cStatus)
				{
					case CStatus.ROOT:
						int prior = result;
						result += combined;
						return prior;
					case CStatus.SECOND:
						secondValue = combined;
						locked = false;

						Monitor.PulseAll(this); // wake up waiting threads

						while (cStatus != CStatus.RESULT) Monitor.Wait(this);
						locked = false;
						Monitor.PulseAll(this);

						cStatus = CStatus.IDLE;
						return result;
					default:
						throw new Exception("unexpected Node state");
				}
			}

		}

		public void distribute(int prior)
		{
			lock (this)
			{
				switch (cStatus)
				{
					case CStatus.FIRST:
						cStatus = CStatus.IDLE;
						locked = false;
						break;
					case CStatus.SECOND:
						result = prior + firstValue;
						cStatus = CStatus.RESULT;
						break;
					default:
						throw new Exception("unexpected Node state");
				}
				Monitor.PulseAll(this);
			}


		}
	}
}
