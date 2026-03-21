using System;
using FifoQueue.Storage;
using Xunit;

namespace FifoQueue.Tests
{
    public class DatabaseQueueTests
    {
        private static IFifoQueue<int> CreateQueue()
            => new FifoQueue<int>(new DatabaseStorageProvider<int>("Data Source=:memory:"));

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
        public void Peek_ReturnsFirstWithoutRemoving()
        {
            using var queue = CreateQueue();
            queue.Enqueue(42);

            Assert.Equal(42, queue.Peek());
            Assert.Equal(42, queue.Peek());
            Assert.Equal(1, queue.Count);
        }

        [Fact]
        public void Count_ReflectsQueueState()
        {
            using var queue = CreateQueue();
            Assert.Equal(0, queue.Count);

            queue.Enqueue(1);
            queue.Enqueue(2);
            Assert.Equal(2, queue.Count);

            queue.Dequeue();
            Assert.Equal(1, queue.Count);
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
        public void WorksWithStringType()
        {
            using var queue = new FifoQueue<string>(
                new DatabaseStorageProvider<string>("Data Source=:memory:"));

            queue.Enqueue("hello");
            queue.Enqueue("world");

            Assert.Equal("hello", queue.Dequeue());
            Assert.Equal("world", queue.Dequeue());
        }

        [Fact]
        public void WorksWithComplexType()
        {
            using var queue = new FifoQueue<TestItem>(
                new DatabaseStorageProvider<TestItem>("Data Source=:memory:"));

            var item1 = new TestItem { Id = 1, Name = "First" };
            var item2 = new TestItem { Id = 2, Name = "Second" };

            queue.Enqueue(item1);
            queue.Enqueue(item2);

            var dequeued = queue.Dequeue();
            Assert.Equal(1, dequeued.Id);
            Assert.Equal("First", dequeued.Name);
        }

        [Fact]
        public void Factory_CreatesDatabaseQueue()
        {
            using var queue = FifoQueueFactory.Create<int>(StorageType.Database, "Data Source=:memory:");

            queue.Enqueue(5);
            queue.Enqueue(10);

            Assert.Equal(5, queue.Dequeue());
            Assert.Equal(10, queue.Dequeue());
        }

        [Fact]
        public void Factory_ThrowsWhenConnectionStringMissing()
        {
            Assert.Throws<ArgumentException>(() =>
                FifoQueueFactory.Create<int>(StorageType.Database));
        }

        public class TestItem
        {
            public int Id { get; set; }
            public string Name { get; set; } = string.Empty;
        }
    }
}
