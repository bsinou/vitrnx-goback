namespace FifoQueue
{
    /// <summary>
    /// Specifies the type of storage backend to use for the FIFO queue.
    /// </summary>
    public enum StorageType
    {
        /// <summary>
        /// In-memory storage using a linked list. Fastest but volatile.
        /// </summary>
        Memory,

        /// <summary>
        /// File-based storage using JSON serialization. Persistent across restarts.
        /// </summary>
        File,

        /// <summary>
        /// SQLite database storage. Persistent and suitable for larger datasets.
        /// </summary>
        Database
    }
}
