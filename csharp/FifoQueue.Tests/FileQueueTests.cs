using System;
using System.IO;
using FifoQueue.Storage;
using Xunit;

namespace FifoQueue.Tests
{
    public class FileQueueTests : IDisposable
    {
        private readonly string _testDir;
        private readonly string _testFilePath;

        public FileQueueTests()
        {
            _testDir = Path.Combine(Path.GetTempPath(), "fifoqueue_tests_" + Guid.NewGuid().ToString("N"));
            Directory.CreateDirectory(_testDir);
            _testFilePath = Path.Combine(_testDir, "queue.json");
        }

        public void Dispose()
        {
            if (Directory.Exists(_testDir))
                Directory.Delete(_testDir, recursive: true);
        }

        private IFifoQueue<int> CreateQueue()
            => new FifoQueue<int>(new FileStorageProvider<int>(_testFilePath));

        [Fact]
        public void Enqueue_Dequeue_ReturnsInFifoOrder()
        {
            using var queue = CreateQueue();
            queue.Enqueue(10);
            queue.Enqueue(20);
            queue.Enqueue(30);

            Assert.Equal(10, queue.Dequeue());
            Assert.Equal(20, queue.Dequeue());
            Assert.Equal(30, queue.Dequeue());
        }

        [Fact]
        public void PersistsAcrossInstances()
        {
            // Write items with first instance
            using (var queue = CreateQueue())
            {
                queue.Enqueue(100);
                queue.Enqueue(200);
            }

            // Read items with second instance
            using (var queue = CreateQueue())
            {
                Assert.Equal(2, queue.Count);
                Assert.Equal(100, queue.Dequeue());
                Assert.Equal(200, queue.Dequeue());
            }
        }

        [Fact]
        public void Clear_RemovesAllItems_AndPersists()
        {
            using (var queue = CreateQueue())
            {
                queue.Enqueue(1);
                queue.Enqueue(2);
                queue.Clear();
            }

            using (var queue = CreateQueue())
            {
                Assert.True(queue.IsEmpty);
            }
        }

        [Fact]
        public void WorksWithStringType()
        {
            var filePath = Path.Combine(_testDir, "string_queue.json");
            using var queue = new FifoQueue<string>(new FileStorageProvider<string>(filePath));

            queue.Enqueue("alpha");
            queue.Enqueue("beta");

            Assert.Equal("alpha", queue.Dequeue());
            Assert.Equal("beta", queue.Dequeue());
        }

        [Fact]
        public void Factory_CreatesFileQueue()
        {
            var filePath = Path.Combine(_testDir, "factory_queue.json");
            using var queue = FifoQueueFactory.Create<int>(StorageType.File, filePath);

            queue.Enqueue(5);
            queue.Enqueue(10);

            Assert.Equal(5, queue.Dequeue());
            Assert.Equal(10, queue.Dequeue());
        }

        [Fact]
        public void Factory_ThrowsWhenPathMissing()
        {
            Assert.Throws<ArgumentException>(() =>
                FifoQueueFactory.Create<int>(StorageType.File));
        }
    }
}
