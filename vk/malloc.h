#ifndef __MALLOC_H__
#define __MALLOC_H__

#include <vulkan/vulkan.h>

extern VkInstance * new_VkInstance();
extern VkInstanceCreateInfo * new_VkInstanceCreateInfo();
extern VkPhysicalDevice * new_VkPhysicalDevice();
extern VkDebugUtilsMessengerCreateInfoEXT * new_VkDebugUtilsMessengerCreateInfoEXT();
extern VkDebugUtilsMessengerEXT * new_VkDebugUtilsMessengerEXT();
extern VkDevice * new_VkDevice();
extern VkPhysicalDeviceFeatures * new_VkPhysicalDeviceFeatures();
extern VkDeviceCreateInfo * new_VkDeviceCreateInfo();
extern VkDeviceQueueCreateInfo * new_VkDeviceQueueCreateInfo();
extern VkQueue * new_VkQueue();
extern VkSurfaceKHR * new_VkSurfaceKHR();
extern VkSwapchainKHR * new_VkSwapchainKHR();
extern VkShaderModule * new_VkShaderModule();
extern VkPipelineLayout * new_VkPipelineLayout();
extern VkRenderPass * new_VkRenderPass();
extern VkCommandPool * new_VkCommandPool();
extern VkSemaphore * new_VkSemaphore();
extern VkFence * new_VkFence();
extern VkBuffer * new_VkBuffer();
extern VkDeviceMemory * new_VkDeviceMemory();

extern VkPipeline * new_VkPipelines(size_t n);
extern VkAttachmentDescription * new_VkAttachmentDescriptions(size_t n);
extern VkAttachmentReference * new_VkAttachmentReferences(size_t n);
extern VkSubpassDescription * new_VkSubpassDescriptions(size_t n);
extern VkViewport * new_VkViewports(size_t n);
extern VkRect2D * new_VkRect2Ds(size_t n);
extern VkPipelineColorBlendAttachmentState * new_VkPipelineColorBlendAttachmentStates(size_t n);
extern VkGraphicsPipelineCreateInfo * new_VkGraphicsPipelineCreateInfos(size_t n);
extern VkFramebuffer * new_VkFramebuffers(size_t n);
extern VkImageView * new_VkImageViews(size_t n);
extern VkCommandBuffer * new_VkCommandBuffers(size_t n);
extern VkClearValue * new_VkClearValues(size_t n);
extern VkSemaphore * new_VkSemaphores(size_t n);
extern VkSubmitInfo * new_VkSubmitInfos(size_t n);
extern VkSubpassDependency * new_VkSubpassDependencys(size_t n);
extern VkSwapchainKHR * new_VkSwapchainKHRs(size_t n);

#endif // #ifndef __MALLOC_H__
