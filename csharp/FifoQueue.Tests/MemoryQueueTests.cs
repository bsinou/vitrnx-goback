using System;
using System.Threading;
using System.Threading.Tasks;
using FifoQueue.Storage;
using Xunit;

namespace FifoQueue.Tests
{
    public class MemoryQueueTests
    {
        private IFifoQueue<int> CreateQueue()
            => new FifoQueue<int>(new MemoryStorageProvider<int>());

        [Fact]
        public void Enqueue_Dequeue_ReturnsInFifoOrder()
        {
            using var queue = CreateQueue();
            queue.Enqueue(1);
            queue.Enqueue(2);
            queue.Enqueue(3);

            Assert.Equal(1, queue.Dequeue());
            Assert.Equal(2, queue.Dequeue());
            Assert.Equal(3, queue.Dequeue());
        }

        [Fact]
        public void Dequeue_EmptyQueue_ThrowsInvalidOperation()
        {
            using var queue = CreateQueue();
            Assert.Throws<InvalidOperationException>(() => queue.Dequeue());
        }

        [Fact]
        public void TryDequeue_EmptyQueue_ReturnsFalse()
        {
            using var queue = CreateQueue();
            Assert.False(queue.TryDequeue(out var item));
            Assert.Equal(default, item);
        }

        [Fact]
        public void TryDequeue_NonEmptyQueue_ReturnsTrueAndItem()
        {
            using var queue = CreateQueue();
            queue.Enqueue(42);
            Assert.True(queue.TryDequeue(out var item));
            Assert.Equal(42, item);
        }

        [Fact]
        public void Peek_ReturnsFirstWithoutRemoving()
        {
            using var queue = CreateQueue();
            queue.Enqueue(10);
            queue.Enqueue(20);

            Assert.Equal(10, queue.Peek());
            Assert.Equal(10, queue.Peek()); // Still the same
            Assert.Equal(2, queue.Count);
        }

        [Fact]
        public void Peek_EmptyQueue_ThrowsInvalidOperation()
        {
            using var queue = CreateQueue();
            Assert.Throws<InvalidOperationException>(() => queue.Peek());
        }

        [Fact]
        public void TryPeek_EmptyQueue_ReturnsFalse()
        {
            using var queue = CreateQueue();
            Assert.False(queue.TryPeek(out _));
        }

        [Fact]
        public void Count_ReflectsQueueState()
        {
            using var queue = CreateQueue();
            Assert.Equal(0, queue.Count);

            queue.Enqueue(1);
            Assert.Equal(1, queue.Count);

            queue.Enqueue(2);
            Assert.Equal(2, queue.Count);

            queue.Dequeue();
            Assert.Equal(1, queue.Count);
        }

        [Fact]
        public void IsEmpty_ReflectsQueueState()
        {
            using var queue = CreateQueue();
            Assert.True(queue.IsEmpty);

            queue.Enqueue(1);
            Assert.False(queue.IsEmpty);

            queue.Dequeue();
            Assert.True(queue.IsEmpty);
        }

        [Fact]
        public void Clear_RemovesAllItems()
        {
            using var queue = CreateQueue();
            queue.Enqueue(1);
            queue.Enqueue(2);
            queue.Enqueue(3);

            queue.Clear();

            Assert.True(queue.IsEmpty);
            Assert.Equal(0, queue.Count);
        }

        [Fact]
        public void Dispose_PreventsSubsequentOperations()
        {
            var queue = CreateQueue();
            queue.Enqueue(1);
            queue.Dispose();

            Assert.Throws<ObjectDisposedException>(() => queue.Enqueue(2));
            Assert.Throws<ObjectDisposedException>(() => queue.Dequeue());
            Assert.Throws<ObjectDisposedException>(() => queue.Peek());
            Assert.Throws<ObjectDisposedException>(() => _ = queue.Count);
        }

        [Fact]
        public void ConcurrentAccess_MaintainsIntegrity()
        {
            using var queue = new FifoQueue<int>(new MemoryStorageProvider<int>());
            const int itemCount = 1000;
            int dequeued = 0;

            // Enqueue from multiple threads
            Parallel.For(0, itemCount, i => queue.Enqueue(i));

            Assert.Equal(itemCount, queue.Count);

            // Dequeue from multiple threads
            Parallel.For(0, itemCount, _ =>
            {
                if (queue.TryDequeue(out _))
                    Interlocked.Increment(ref dequeued);
            });

            Assert.Equal(itemCount, dequeued);
            Assert.True(queue.IsEmpty);
        }

        [Fact]
        public void Factory_CreatesMemoryQueue()
        {
            using var queue = FifoQueueFactory.Create<string>(StorageType.Memory);
            queue.Enqueue("hello");
            queue.Enqueue("world");

            Assert.Equal("hello", queue.Dequeue());
            Assert.Equal("world", queue.Dequeue());
        }
    }
}
