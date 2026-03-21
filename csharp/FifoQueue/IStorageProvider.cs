using System;

namespace FifoQueue
{
    /// <summary>
    /// Abstraction for the underlying storage mechanism used by the FIFO queue.
    /// Implementations must handle their own data persistence but are not required to be thread-safe;
    /// thread safety is ensured by <see cref="FifoQueue{T}"/>.
    /// </summary>
    /// <typeparam name="T">The type of elements stored.</typeparam>
    public interface IStorageProvider<T> : IDisposable
    {
        /// <summary>
        /// Appends an item to the end of the storage.
        /// </summary>
        void Append(T item);

        /// <summary>
        /// Removes and returns the first item from the storage.
        /// </summary>
        /// <returns>The first item.</returns>
        /// <exception cref="InvalidOperationException">Thrown when storage is empty.</exception>
        T RemoveFirst();

        /// <summary>
        /// Attempts to remove and return the first item from the storage.
        /// </summary>
        /// <param name="item">The removed item, or the default value if storage is empty.</param>
        /// <returns><c>true</c> if an item was removed; otherwise, <c>false</c>.</returns>
        bool TryRemoveFirst(out T? item);

        /// <summary>
        /// Returns the first item without removing it.
        /// </summary>
        /// <returns>The first item.</returns>
        /// <exception cref="InvalidOperationException">Thrown when storage is empty.</exception>
        T PeekFirst();

        /// <summary>
        /// Attempts to return the first item without removing it.
        /// </summary>
        /// <param name="item">The first item, or the default value if storage is empty.</param>
        /// <returns><c>true</c> if an item was found; otherwise, <c>false</c>.</returns>
        bool TryPeekFirst(out T? item);

        /// <summary>
        /// Gets the number of items in storage.
        /// </summary>
        int Count { get; }

        /// <summary>
        /// Removes all items from storage.
        /// </summary>
        void Clear();
    }
}
