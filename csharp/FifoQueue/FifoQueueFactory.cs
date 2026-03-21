using System;
using FifoQueue.Storage;

namespace FifoQueue
{
    /// <summary>
    /// Factory for creating FIFO queues with a specific storage backend.
    /// </summary>
    public static class FifoQueueFactory
    {
        /// <summary>
        /// Creates a new FIFO queue using the specified storage type.
        /// </summary>
        /// <typeparam name="T">The type of elements in the queue.</typeparam>
        /// <param name="storageType">The type of storage backend.</param>
        /// <param name="connectionStringOrPath">
        /// For <see cref="StorageType.File"/>: the file path for persistence.
        /// For <see cref="StorageType.Database"/>: the SQLite connection string (e.g., "Data Source=queue.db").
        /// For <see cref="StorageType.Memory"/>: this parameter is ignored and can be null.
        /// </param>
        /// <returns>A new <see cref="IFifoQueue{T}"/> instance.</returns>
        /// <exception cref="ArgumentException">
        /// Thrown when a required <paramref name="connectionStringOrPath"/> is missing for File or Database storage.
        /// </exception>
        /// <exception cref="ArgumentOutOfRangeException">
        /// Thrown when an unsupported <paramref name="storageType"/> is provided.
        /// </exception>
        public static IFifoQueue<T> Create<T>(StorageType storageType, string? connectionStringOrPath = null)
        {
            IStorageProvider<T> provider = storageType switch
            {
                StorageType.Memory => new MemoryStorageProvider<T>(),
                StorageType.File => new FileStorageProvider<T>(
                    connectionStringOrPath ?? throw new ArgumentException(
                        "A file path is required for File storage.", nameof(connectionStringOrPath))),
                StorageType.Database => new DatabaseStorageProvider<T>(
                    connectionStringOrPath ?? throw new ArgumentException(
                        "A connection string is required for Database storage.", nameof(connectionStringOrPath))),
                _ => throw new ArgumentOutOfRangeException(nameof(storageType), storageType, "Unsupported storage type.")
            };

            return new FifoQueue<T>(provider);
        }
    }
}
