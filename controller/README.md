k8s 怎么编写自定义的 controller 
三种方式：
- 使用 informer
- 使用 Controller runtime
- 使用 Kubebuilder

下面是采用 Informer，Controller runtime 和 Kubebuilder 来编写 Controller 的区别：

直接使用 Informer：直接使用 Informer 编写 Controller 需要编写更多的代码，因为我们需要在代码处理更多的底层细节，例如如何在集群中监视资源，以及如何处理资源变化的通知。但是，使用 Informer 也可以更加自定义和灵活，因为我们可以更细粒度地控制 Controller 的行为。

Controller runtime：Controller runtime 是基于 Informer 实现的，在 Informer 之上为 Controller 编写提供了高级别的抽象和帮助类，包括 Leader Election、Event Handling 和 Reconcile Loop 等等。使用 Controller runtime，可以更容易地编写和测试 Controller，因为它已经处理了许多底层的细节。

Kubebuilder：和 Informer 及 Controller runtime 不同，Kubebuilder 并不是一个代码库，而是一个开发框架。Kubebuilder 底层使用了 controller-runtime。Kubebuilder 提供了 CRD 生成器和代码生成器等工具，可以帮助开发者自动生成一些重复性的代码和资源定义，提高开发效率。同时，Kubebuilder 还可以生成 Webhooks，以用于验证自定义资源。