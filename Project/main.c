// file:../Compiler-build/../Compiler-build/hello_world.ref

#include <memory_manager.h>
void __initLiteralData()
{
 initAllocator(1024 * 1024 * 1024);

 initHeaps(2);
} // __initLiteralData()

int main()
{
 __initLiteralData();
 return 0;
}
