package vk

// #include "malloc.h"
//
// #include <stdlib.h>
//
// VkInstance * new_VkInstance() {
//    return calloc(1, sizeof(VkInstance));
// }
//
// VkInstanceCreateInfo * new_VkInstanceCreateInfo() {
//    return calloc(1, sizeof(VkInstanceCreateInfo));
// }
//
// VkPhysicalDevice * new_VkPhysicalDevice() {
//    return calloc(1, sizeof(VkPhysicalDevice));
// }
//
// VkDebugUtilsMessengerCreateInfoEXT * new_VkDebugUtilsMessengerCreateInfoEXT() {
//    return calloc(1, sizeof(VkDebugUtilsMessengerCreateInfoEXT));
// }
//
// VkDebugUtilsMessengerEXT * new_VkDebugUtilsMessengerEXT() {
//    return calloc(1, sizeof(VkDebugUtilsMessengerEXT));
// }
//
// VkDevice * new_VkDevice() {
//    return calloc(1, sizeof(VkDevice));
// }
//
// VkPhysicalDeviceFeatures * new_VkPhysicalDeviceFeatures() {
//    return calloc(1, sizeof(VkPhysicalDeviceFeatures));
// }
//
// VkDeviceCreateInfo * new_VkDeviceCreateInfo() {
//    return calloc(1, sizeof(VkDeviceCreateInfo));
// }
//
// VkDeviceQueueCreateInfo * new_VkDeviceQueueCreateInfo() {
//    return calloc(1, sizeof(VkDeviceQueueCreateInfo));
// }
//
// VkQueue * new_VkQueue() {
//    return calloc(1, sizeof(VkQueue));
// }
//
// VkSurfaceKHR * new_VkSurfaceKHR() {
//    return calloc(1, sizeof(VkSurfaceKHR));
// }
//
// VkSwapchainKHR * new_VkSwapchainKHR() {
//    return calloc(1, sizeof(VkSwapchainKHR));
// }
//
// VkShaderModule * new_VkShaderModule() {
//    return calloc(1, sizeof(VkShaderModule));
// }
//
// VkPipelineLayout * new_VkPipelineLayout() {
//    return calloc(1, sizeof(VkPipelineLayout));
// }
//
// VkRenderPass * new_VkRenderPass() {
//    return calloc(1, sizeof(VkRenderPass));
// }
//
// VkCommandPool * new_VkCommandPool() {
//    return calloc(1, sizeof(VkCommandPool));
// }
//
// VkSemaphore * new_VkSemaphore() {
//    return calloc(1, sizeof(VkSemaphore));
// }
//
// VkFence * new_VkFence() {
//    return calloc(1, sizeof(VkFence));
// }
//
// VkBuffer * new_VkBuffer() {
//    return calloc(1, sizeof(VkBuffer));
// }
//
// VkDeviceMemory * new_VkDeviceMemory() {
//    return calloc(1, sizeof(VkDeviceMemory));
// }
//
//
//
// VkPipeline * new_VkPipelines(size_t n) {
//    return calloc(n, sizeof(VkPipeline));
// }
//
// VkAttachmentDescription * new_VkAttachmentDescriptions(size_t n) {
//    return calloc(n, sizeof(VkAttachmentDescription));
// }
//
// VkAttachmentReference * new_VkAttachmentReferences(size_t n) {
//    return calloc(n, sizeof(VkAttachmentReference));
// }
//
// VkSubpassDescription * new_VkSubpassDescriptions(size_t n) {
//    return calloc(n, sizeof(VkSubpassDescription));
// }
//
// VkViewport * new_VkViewports(size_t n) {
//    return calloc(n, sizeof(VkViewport));
// }
//
// VkRect2D * new_VkRect2Ds(size_t n) {
//    return calloc(n, sizeof(VkRect2D));
// }
//
// VkPipelineColorBlendAttachmentState * new_VkPipelineColorBlendAttachmentStates(size_t n) {
//    return calloc(n, sizeof(VkPipelineColorBlendAttachmentState));
// }
//
// VkGraphicsPipelineCreateInfo * new_VkGraphicsPipelineCreateInfos(size_t n) {
//    return calloc(n, sizeof(VkGraphicsPipelineCreateInfo));
// }
//
// VkFramebuffer * new_VkFramebuffers(size_t n) {
//    return calloc(n, sizeof(VkFramebuffer));
// }
//
// VkImageView * new_VkImageViews(size_t n) {
//    return calloc(n, sizeof(VkImageView));
// }
//
// VkCommandBuffer * new_VkCommandBuffers(size_t n) {
//    return calloc(n, sizeof(VkCommandBuffer));
// }
//
// VkClearValue * new_VkClearValues(size_t n) {
//    return calloc(n, sizeof(VkClearValue));
// }
//
// VkSemaphore * new_VkSemaphores(size_t n) {
//    return calloc(n, sizeof(VkSemaphore));
// }
//
// VkSubmitInfo * new_VkSubmitInfos(size_t n) {
//    return calloc(n, sizeof(VkSubmitInfo));
// }
//
// VkSubpassDependency * new_VkSubpassDependencys(size_t n) {
//    return calloc(n, sizeof(VkSubpassDependency));
// }
//
// VkSwapchainKHR * new_VkSwapchainKHRs(size_t n) {
//    return calloc(n, sizeof(VkSwapchainKHR));
// }
import "C"
