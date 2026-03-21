using System;

namespace FifoQueue
{
    /// <summary>
    /// Represents a thread-safe FIFO (First-In, First-Out) queue with pluggable storage.
    /// </summary>
    /// <typeparam name="T">The type of elements in the queue.</typeparam>
    public interface IFifoQueue<T> : IDisposable
    {
        /// <summary>
        /// Adds an item to the end of the queue.
        /// </summary>
        void Enqueue(T item);

        /// <summary>
        /// Removes and returns the item at the front of the queue.
        /// </summary>
        /// <returns>The item at the front of the queue.</returns>
        /// <exception cref="InvalidOperationException">Thrown when the queue is empty.</exception>
        T Dequeue();

        /// <summary>
        /// Attempts to remove and return the item at the front of the queue.
        /// </summary>
        /// <param name="item">The item removed from the front of the queue, or the default value if the queue is empty.</param>
        /// <returns><c>true</c> if an item was removed; otherwise, <c>false</c>.</returns>
        bool TryDequeue(out T? item);

        /// <summary>
        /// Returns the item at the front of the queue without removing it.
        /// </summary>
        /// <returns>The item at the front of the queue.</returns>
        /// <exception cref="InvalidOperationException">Thrown when the queue is empty.</exception>
        T Peek();

        /// <summary>
        /// Attempts to return the item at the front of the queue without removing it.
        /// </summary>
        /// <param name="item">The item at the front of the queue, or the default value if the queue is empty.</param>
        /// <returns><c>true</c> if an item was found; otherwise, <c>false</c>.</returns>
        bool TryPeek(out T? item);

        /// <summary>
        /// Gets the number of items currently in the queue.
        /// </summary>
        int Count { get; }

        /// <summary>
        /// Gets a value indicating whether the queue is empty.
        /// </summary>
        bool IsEmpty { get; }

        /// <summary>
        /// Removes all items from the queue.
        /// </summary>
        void Clear();
    }
}
